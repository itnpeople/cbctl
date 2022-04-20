package tumblebug

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type MCISOptions struct {
	app.Output
	RootUrl   string
	Namespace string
	Name      string
	HTTP      *resty.Request
}

// returns initialized Options
func NewMCISOptions(output app.Output) *MCISOptions {
	return &MCISOptions{
		Output: output,
	}
}

// completes all the required options
func (o *MCISOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, app.Config.GetCurrentContext().Urls.Tumblebug)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid Tumblebug ROOT URL (%s)", o.RootUrl)
	}
	o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Invalid namespace")
	}
	o.HTTP = resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
	return nil
}

// validates the provided options
func (o *MCISOptions) Validate(args []string) error {
	if len(args) > 0 {
		o.Name = utils.NVL(o.Name, args[0])
	}
	if o.Name == "" {
		return fmt.Errorf("Invalid name")
	}
	return nil
}

// returns a cobra command
func NewCmdMCIS(output app.Output) *cobra.Command {
	o := NewMCISOptions(output)
	cmds := &cobra.Command{
		Use:   "mcis",
		Short: "Cloud Infra.",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "Tumblebug ROOT URL (http://localhost:1323/tumblebug)")
	cmds.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "cloud-barista namespace for cluster list")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of MCIS")

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get MCISs.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := o.HTTP.Get(fmt.Sprintf("%s/ns/%s/mcis", o.RootUrl, o.Namespace)); err != nil {
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
		Short: "Get a MCIS.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate(args))
			app.ValidateError(func() error {
				if resp, err := o.HTTP.Get(fmt.Sprintf("%s/ns/%s/mcis/%s", o.RootUrl, o.Namespace, o.Name)); err != nil {
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
		Short: "Delete a MCIS.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate(args))
			app.ValidateError(func() error {
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/mcis/%s?action=terminate", o.RootUrl, o.Namespace, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/mcis/%s?action=refine", o.RootUrl, o.Namespace, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}

				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/mcis/%s", o.RootUrl, o.Namespace, o.Name)); err != nil {
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
		Use:   "clean",
		Short: "clean-up MCISs.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {

				// mcis
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/mcis", o.RootUrl, o.Namespace)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				// vpc
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/resources/vNet", o.RootUrl, o.Namespace)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				// securityGroup
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/resources/securityGroup", o.RootUrl, o.Namespace)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				// sshKey
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/resources/sshKey", o.RootUrl, o.Namespace)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				// image
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/resources/image", o.RootUrl, o.Namespace)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				// spec
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s/resources/spec", o.RootUrl, o.Namespace)); err != nil {
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
