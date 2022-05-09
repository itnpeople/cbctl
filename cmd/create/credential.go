package create

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type CredentialOptions struct {
	*app.Options
	CSP            string
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

// validates the provided options
func (o *CredentialOptions) Validate() error {
	o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}

	if o.Filename != "" {
		return nil
	}

	if o.Name == "" {
		return fmt.Errorf("Name is required.")
	}
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
func NewCommandCredential(options *app.Options) *cobra.Command {
	o := &CredentialOptions{
		Options: options,
	}

	// create
	cmd := &cobra.Command{
		Use:                   "credential (NAME | --name NAME | -f FILENAME) [options]",
		Short:                 "Create a cloud credential",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, o.Validate())
			app.ValidateError(c, func() error {
				if out, err := app.GetBody(o, `{
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
					]
				}`); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/credential", app.Config.GetCurrentContext().Urls.Spider)
					if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
						return err
					} else {
						o.WriteBody(resp.Body())
					}
				}
				return nil
			}())
		},
	}

	cmd.Flags().StringVar(&o.CSP, "csp", "", "Cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack, cloudit)")
	cmd.Flags().StringVar(&o.ClientID, "secret-id", "", "Key id (aws, azure, alibaba, tencent)")      // AWS, Azure, Alibaba, Tencet
	cmd.Flags().StringVar(&o.ClientSecret, "secret", "", "Key secret (aws, azure, alibaba, tencent)") // AWS, Azure, Alibaba, Tencet
	cmd.Flags().StringVar(&o.SubscriptionId, "subscription", "", "Subscription id (azure)")           // Azure (additional)
	cmd.Flags().StringVar(&o.TenantId, "tenant", "", "Tenant id (azure, cloudit)")                    // Azure, Cloudit (additional)
	cmd.Flags().StringVar(&o.ClientEmail, "client-email", "", "Client email (gcp)")                   // GCP
	cmd.Flags().StringVar(&o.ProjectID, "project-id", "", "Project id (gcp, openstack)")              // GCP, openstack
	cmd.Flags().StringVar(&o.PrivateKey, "private-key", "", "Private-key (gcp)")                      // GCP
	cmd.Flags().StringVar(&o.ApiKey, "api-key", "", "Api-key (ibm)")                                  // IBM
	cmd.Flags().StringVar(&o.Endpoint, "endpoint", "", "Identity Endpoint (openstack, cloudit)")      // Openstack, Cloudit
	cmd.Flags().StringVar(&o.Username, "username", "", "Username (openstack, cloudit)")               // Openstack, Cloudit
	cmd.Flags().StringVar(&o.Password, "password", "", "Password (openstack, cloudit)")               // Openstack, Cloudit
	cmd.Flags().StringVar(&o.DomainName, "domain", "", "Domain Name (openstack)")                     // Openstack
	cmd.Flags().StringVar(&o.AutoToken, "token", "", "Auth Token (cloudit)")                          // Cloudit

	return cmd

}
