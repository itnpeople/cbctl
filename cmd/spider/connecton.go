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
type ConnectionOptions struct {
	app.ConfigContext
	app.Output
	RootUrl    string
	Name       string
	CSP        string
	Credential string
	Region     string
}

// returns initialized Options
func NewConnectionOptions(ctx app.ConfigContext, output app.Output) *ConnectionOptions {
	return &ConnectionOptions{
		ConfigContext: ctx,
		Output:        output,
	}
}

// completes all the required options
func (o *ConnectionOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, o.ConfigContext.Urls.Spider)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request roo-url flag (%s)", o.RootUrl)
	}
	return nil
}

// validates the provided options
func (o *ConnectionOptions) Validate() error {
	if o.CSP == "" {
		return fmt.Errorf("Invalid csp flag")
	}
	if o.Name == "" {
		return fmt.Errorf("Invalid name flag")
	}
	if o.Credential == "" {
		return fmt.Errorf("Invalid credential flag")
	}
	if o.Region == "" {
		return fmt.Errorf("Invalid region flag")
	}
	return nil
}

// returns a cobra command
func NewCmdConnection(ctx app.ConfigContext, output app.Output) *cobra.Command {
	o := NewConnectionOptions(ctx, output)
	cmds := &cobra.Command{
		Use:   "connection",
		Short: "Cloud connection info.",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "spider root url (http://localhost:1024/spider)")
	cmds.PersistentFlags().StringVar(&o.CSP, "csp", "", "cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack)")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of connection info.")

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a cloud connection info.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if out, err := utils.ToTemplateBytes(tplConnection, o); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/connectionconfig", o.RootUrl)
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
	cmdC.Flags().StringVar(&o.Region, "region", "", "region")
	cmdC.Flags().StringVar(&o.Credential, "credential", "", "credential")
	cmds.AddCommand(cmdC)

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get cloud connection infos.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/connectionconfig", o.RootUrl)); err != nil {
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
		Short: "Get a cloud connection infos.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/connectionconfig/%s", o.RootUrl, o.Name)); err != nil {
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
		Short: "Delete a cloud connection info.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/connectionconfig/%s", o.RootUrl, o.Name)); err != nil {
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
	tplConnection = `{
"ConfigName"     : "{{ .Name }}",
"ProviderName"   : "{{ .CSP | ToUpper }}", 
"DriverName"     : "{{ .CSP }}-driver-v1.0", 
"CredentialName" : "{{ .Credential }}", 
"RegionName"     : "{{ .Region }}"
}`
)
