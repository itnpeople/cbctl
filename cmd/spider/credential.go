package spider

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type CredentialOptions struct {
	app.ConfigContext
	app.Output
	RootUrl        string
	CSP            string
	Name           string
	ClientID       string
	ClientSecret   string
	ClientEmail    string
	ProjectID      string
	PrivateKey     string
	TenantId       string
	SubscriptionId string
	ApiKey         string
}

// returns initialized Options
func NewCredentialOptions(ctx app.ConfigContext, output app.Output) *CredentialOptions {
	return &CredentialOptions{
		ConfigContext: ctx,
		Output:        output,
	}
}

// completes all the required options
func (o *CredentialOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, o.ConfigContext.Urls.Spider)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request roo-url flag (%s)", o.RootUrl)
	}
	return nil
}

// validates the provided options
func (o *CredentialOptions) Validate() error {
	switch o.CSP {
	case "aws", "alibaba", "tencet":
		if o.ClientID == "" || o.ClientSecret == "" {
			return fmt.Errorf("Invalid credential flag (csp=%s, key=%s, secret=%s)", o.CSP, o.ClientID, o.ClientSecret)
		}
		break
	case "gcp":
		if o.ClientEmail == "" || o.ProjectID == "" || o.PrivateKey == "" {
			return fmt.Errorf("Invalid credential flag (csp=%s, email=%s, project=%s, private-key=%s)", o.CSP, o.ClientEmail, o.ProjectID, o.PrivateKey)
		}
		break
	case "azure":
		if o.ClientID == "" || o.ClientSecret == "" || o.TenantId == "" || o.SubscriptionId == "" {
			return fmt.Errorf("Invalid credential flag (csp=%s, id=%s, secret=%s, tenant=%s, subscription=%s)", o.CSP, o.ClientID, o.ClientSecret, o.TenantId, o.SubscriptionId)
		}
		break
	case "ibm":
		if o.ApiKey == "" {
			return fmt.Errorf("Invalid credential flag (csp=%s, api-key=%s)", o.CSP, o.ApiKey)
		}
		break
	default:
		return fmt.Errorf("Not supported CSP (csp=%s)", o.CSP)
	}
	return nil
}

// returns a cobra command
func NewCmdCredential(ctx app.ConfigContext, output app.Output) *cobra.Command {
	o := NewCredentialOptions(ctx, output)
	cmds := &cobra.Command{
		Use:   "credential",
		Short: "Cloud credential",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "spider root url (http://localhost:1024/spider)")
	cmds.PersistentFlags().StringVar(&o.CSP, "csp", "", "cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack)")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of credential")

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a cloud credential",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if out, err := utils.ToTemplateBytes(tplCredential, o); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/credential", o.RootUrl)
					if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
						return err
					} else {
						o.Output.Write(resp.Body())
					}
				}
				return nil
			}())
		},
	}
	cmdC.Flags().StringVar(&o.ClientID, "secret-id", "", "key id (aws, azure, alibaba, tencent)")      // AWS, Azure, Alibaba, Tencet
	cmdC.Flags().StringVar(&o.ClientSecret, "secret", "", "key secret (aws, azure, alibaba, tencent)") // AWS, Azure, Alibaba, Tencet
	cmdC.Flags().StringVar(&o.SubscriptionId, "subscription-id", "", "Azure subscription id")          // Azure (additional)
	cmdC.Flags().StringVar(&o.TenantId, "tenant", "", "Azure tenant id")                               // Azure (additional)
	cmdC.Flags().StringVar(&o.ClientEmail, "client-email", "", "Google Cloud client email")            // GCP
	cmdC.Flags().StringVar(&o.ProjectID, "project-id", "", "Google Cloud project id")                  // GCP
	cmdC.Flags().StringVar(&o.PrivateKey, "private-key", "", "Google Cloud private-key")               // GCP
	cmdC.Flags().StringVar(&o.ApiKey, "api-key", "", "IBM api-key")                                    // IBM
	cmds.AddCommand(cmdC)

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get a cloud credentials",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/credential", o.RootUrl)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// get
	cmds.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a cloud credential",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/credential/%s", o.RootUrl, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// delete
	cmds.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a cloud credential",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/credential/%s", o.RootUrl, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	return cmds
}

const (
	tplCredential = `{
"CredentialName"   : "{{ .Name }}",
"ProviderName"     : "{{ .CSP | ToUpper }}",
"KeyValueInfoList" : [
	{"Key" : "ClientId",       "Value" : "{{ .ClientID }}"},
	{"Key" : "ClientSecret",   "Value" : "{{ .ClientSecret }}"},
	{"Key" : "ClientEmail",    "Value" : "{{ .ClientEmail }}"},
	{"Key" : "ProjectID",      "Value" : "{{ .ProjectID }}"},
	{"Key" : "PrivateKey",     "Value" : "{{ .PrivateKey }}"},
	{"Key" : "TenantId",       "Value" : "{{ .TenantId }}"},
	{"Key" : "SubscriptionId", "Value" : "{{ .SubscriptionId }}"},
	{"Key" : "ApiKey",         "Value" : "{{ .ApiKey }}"}
]}`
)
