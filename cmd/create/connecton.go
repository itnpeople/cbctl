package create

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
)

// a struct to support command
type ConnectionOptions struct {
	*app.Options
	CSP        string
	Credential string
	Region     string
}

// returns a cobra command
func NewCommandConnection(options *app.Options) *cobra.Command {
	o := &ConnectionOptions{
		Options: options,
	}

	// create
	cmd := &cobra.Command{
		Use:                   "connection (NAME | --name NAME | -f FILENAME) [options]",
		Short:                 "Create a cloud connection info.",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Filename == "" {
					if o.Name == "" {
						return fmt.Errorf("Name is required.")
					}
					if o.CSP == "" {
						return fmt.Errorf("CSP is required.")
					}
					if o.Credential == "" {
						return fmt.Errorf("Credential name is required.")
					}
					if o.Region == "" {
						return fmt.Errorf("Region name is required.")
					}
				}
				if out, err := app.GetBody(o, `{
					"ConfigName"     : "{{ .Name }}",
					"ProviderName"   : "{{ .CSP | ToUpper }}", 
					"DriverName"     : "{{ .CSP }}-driver-v1.0", 
					"CredentialName" : "{{ .Credential }}", 
					"RegionName"     : "{{ .Region }}"
				}`); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/connectionconfig", app.Config.GetCurrentContext().Urls.Spider)
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
	cmd.Flags().StringVar(&o.Region, "region", "", "Region name")
	cmd.Flags().StringVar(&o.Credential, "credential", "", "Credential name")

	return cmd
}
