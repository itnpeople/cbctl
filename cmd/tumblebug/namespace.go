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
type NamespaceOptions struct {
	app.Output
	RootUrl     string
	HTTP        *resty.Request
	Name        string
	Description string
}

// returns initialized Options
func NewNamespaceOptions(output app.Output) *NamespaceOptions {
	return &NamespaceOptions{
		Output: output,
	}
}

// completes all the required options
func (o *NamespaceOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, app.Config.GetCurrentContext().Urls.Tumblebug)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid Tumblebug ROOT URL (%s)", o.RootUrl)
	}
	o.HTTP = resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
	return nil
}

// returns a cobra command
func NewCmdNamespace(output app.Output) *cobra.Command {
	o := NewNamespaceOptions(output)
	cmds := &cobra.Command{
		Use:   "namespace",
		Short: "Cloud-barista Namespace",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "Tumblebug ROOT URL (http://localhost:1323/tumblebug)")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of MCIS")

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get namespaces.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := o.HTTP.Get(fmt.Sprintf("%s/ns", o.RootUrl)); err != nil {
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
		Short: "Get a namespace.",
		Args:  app.ValidCommandArgs(0, &o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := o.HTTP.Get(fmt.Sprintf("%s/ns/%s", o.RootUrl, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a namespace.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				var body = fmt.Sprintf("{\"name\" : \"%s\", \"description\" : \"%s\"}", o.Name, o.Description)
				if resp, err := o.HTTP.SetHeader("content-type", "application/json").SetBody(body).Post(fmt.Sprintf("%s/ns", o.RootUrl)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	}
	cmdC.Flags().StringVar(&o.Description, "description", "", "description")
	cmds.AddCommand(cmdC)

	// delete
	cmds.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a namespace.",
		Args:  app.ValidCommandArgs(0, &o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := o.HTTP.Delete(fmt.Sprintf("%s/ns/%s", o.RootUrl, o.Name)); err != nil {
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
	tplNamespace = `{
	"name"        : "{{NAME}}",
	"description" : ""
}}`
)
