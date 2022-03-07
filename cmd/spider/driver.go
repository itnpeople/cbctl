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
type DriverOptions struct {
	app.Output
	app.ConfigContext
	RootUrl string
	CSP     string
}

// returns initialized Options
func NewDriverOptions(ctx app.ConfigContext, output app.Output) *DriverOptions {
	return &DriverOptions{
		ConfigContext: ctx,
		Output:        output,
	}
}

// completes all the required options
func (o *DriverOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, o.ConfigContext.Urls.Spider)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request roo-url flag (%s)", o.RootUrl)
	}
	return nil
}

// validates the provided options
func (o *DriverOptions) Validate() error {
	if o.CSP == "" {
		return fmt.Errorf("Invalid csp flag")
	}
	return nil
}

// returns a cobra command
func NewCmdDriver(ctx app.ConfigContext, output app.Output) *cobra.Command {
	o := NewDriverOptions(ctx, output)
	cmds := &cobra.Command{
		Use:   "driver",
		Short: "Cloud driver",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "spider root url (http://localhost:1024/spider)")
	cmds.PersistentFlags().StringVar(&o.CSP, "csp", "", "cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack)")

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get cloud drivers",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/driver", o.RootUrl)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// create
	cmds.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a cloud driver",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if out, err := utils.ToTemplateBytes(tplDriver, o); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/driver", o.RootUrl)
					if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
						return err
					} else {
						o.Output.Write(resp.Body())
					}
				}
				return nil
			}())
		},
	})

	// get
	cmds.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a cloud driver",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/driver/%s-driver-v1.0", o.RootUrl, o.CSP)); err != nil {
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
		Short: "Delete a cloud driver",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/driver/%s-driver-v1.0", o.RootUrl, o.CSP)); err != nil {
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
	tplDriver = `{
"DriverName"        : "{{ .CSP }}-driver-v1.0",
"ProviderName"      : "{{ .CSP | ToUpper }}",
"DriverLibFileName" : "{{ .CSP }}-driver-v1.0.so"
}`
)
