package create

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
)

// a struct to support command
type DriverOptions struct {
	*app.Options
	CSP string
}

// returns a cobra command
func NewCommandDriver(options *app.Options) *cobra.Command {
	o := &DriverOptions{
		Options: options,
	}

	// create
	cmd := &cobra.Command{
		Use:                   "driver (CSP | --csp CSP | -f FILENAME) [options]",
		Short:                 "Create a cloud driver",
		Args:                  app.BindCommandArgs(&o.CSP),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Filename == "" && o.CSP == "" {
					return fmt.Errorf("CSP is required.")
				}
				if out, err := app.GetBody(o, `{
					"DriverName"        : "{{ .CSP }}-driver-v1.0",
					"ProviderName"      : "{{ .CSP | ToUpper }}",
					"DriverLibFileName" : "{{ .CSP }}-driver-v1.0.so"
				}`); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/driver", app.Config.GetCurrentContext().Urls.Spider)
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

	return cmd

}
