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
	Endpoint       string
	Username       string
	Password       string
	DomainName     string
	AutoToken      string
}

// returns initialized Options
func NewCredentialOptions(output app.Output) *CredentialOptions {
	return &CredentialOptions{
		Output: output,
	}
}

// completes all the required options
func (o *CredentialOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, app.Config.GetCurrentContext().Urls.Spider)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request root-url flag (%s)", o.RootUrl)
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
	case "openstack":
		if o.Endpoint == "" || o.Username == "" || o.Password == "" || o.DomainName == "" || o.ProjectID == "" {
			return fmt.Errorf("Invalid credential flag (csp=%s, endpoint=%s, username=%s, password=%s, domain=%s, project-id=%s)", o.CSP, o.Endpoint, o.Username, o.Password, o.DomainName, o.ProjectID)
		}
		break
	case "cloudit":
		if o.Endpoint == "" || o.Username == "" || o.Password == "" || o.AutoToken == "" || o.TenantId == "" {
			return fmt.Errorf("Invalid credential flag (csp=%s, endpoint=%s, username=%s, password=%s, token=%s, tenent=%s)", o.CSP, o.Endpoint, o.Username, o.Password, o.AutoToken, o.TenantId)
		}
		break
	default:
		return fmt.Errorf("Not supported CSP (csp=%s)", o.CSP)
	}
	return nil
}

// returns a cobra command
func NewCmdCredential(output app.Output) *cobra.Command {
	o := NewCredentialOptions(output)
	cmds := &cobra.Command{
		Use:   "credential",
		Short: "Cloud credential",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "spider root url (http://localhost:1024/spider)")
	cmds.PersistentFlags().StringVar(&o.CSP, "csp", "", "cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack, cloudit)")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of credential")

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a cloud credential",
		Args:  app.ValidCommandArgs(0, &o.Name),
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
	cmdC.Flags().StringVar(&o.ClientID, "secret-id", "", "Key id (aws, azure, alibaba, tencent)")      // AWS, Azure, Alibaba, Tencet
	cmdC.Flags().StringVar(&o.ClientSecret, "secret", "", "Key secret (aws, azure, alibaba, tencent)") // AWS, Azure, Alibaba, Tencet
	cmdC.Flags().StringVar(&o.SubscriptionId, "subscription", "", "Subscription id (azure)")           // Azure (additional)
	cmdC.Flags().StringVar(&o.TenantId, "tenant", "", "Tenant id (azure, cloudit)")                    // Azure, Cloudit (additional)
	cmdC.Flags().StringVar(&o.ClientEmail, "client-email", "", "Client email (gcp)")                   // GCP
	cmdC.Flags().StringVar(&o.ProjectID, "project-id", "", "Project id (gcp, openstack)")              // GCP, openstack
	cmdC.Flags().StringVar(&o.PrivateKey, "private-key", "", "Private-key (gcp)")                      // GCP
	cmdC.Flags().StringVar(&o.ApiKey, "api-key", "", "Api-key (ibm)")                                  // IBM
	cmdC.Flags().StringVar(&o.Endpoint, "endpoint", "", "Identity Endpoint (openstack, cloudit)")      // Openstack, Cloudit
	cmdC.Flags().StringVar(&o.Username, "username", "", "Username (openstack, cloudit)")               // Openstack, Cloudit
	cmdC.Flags().StringVar(&o.Password, "password", "", "Password (openstack, cloudit)")               // Openstack, Cloudit
	cmdC.Flags().StringVar(&o.DomainName, "domain", "", "Domain Name (openstack)")                     // Openstack
	cmdC.Flags().StringVar(&o.AutoToken, "token", "", "Auth Token (cloudit)")                          // Cloudit

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
		Args:  app.ValidCommandArgs(0, &o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
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
		Args:  app.ValidCommandArgs(0, &o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
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
	{"Key" : "ClientId",         "Value" : "{{ .ClientID }}"},
	{"Key" : "ClientSecret",     "Value" : "{{ .ClientSecret }}"},
	{"Key" : "ClientEmail",      "Value" : "{{ .ClientEmail }}"},
	{"Key" : "ProjectID",        "Value" : "{{ .ProjectID }}"},
	{"Key" : "PrivateKey",       "Value" : "{{ .PrivateKey }}"},
	{"Key" : "TenantId",         "Value" : "{{ .TenantId }}"},
	{"Key" : "SubscriptionId",   "Value" : "{{ .SubscriptionId }}"},
	{"Key" : "ApiKey",           "Value" : "{{ .ApiKey }}"},
	{"Key" : "IdentityEndpoint", "Value" : "{{ .Endpoint }}"},
	{"Key" : "Username",         "Value" : "{{ .Username }}"},
	{"Key" : "Password",         "Value" : "{{ .Password }}"},
	{"Key" : "DomainName",       "Value" : "{{ .DomainName }}"},
	{"Key" : "AuthToken",        "Value" : "{{ .AutoToken }}"}
]}`
)
