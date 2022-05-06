package create

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
)

// a struct to support command
type NamespaceOptions struct {
	*app.Options
	Description string
}

// returns a cobra command
func NewCommandNamespace(options *app.Options) *cobra.Command {
	o := &NamespaceOptions{
		Options: options,
	}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a namespace.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if out, err := app.GetBody(o, `{
					"name"        : "{{NAME}}",
					"description" : "{{Description}}"
				}`); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/ns", app.Config.GetCurrentContext().Urls.Tumblebug)
					if resp, err := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default").SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
						return err
					} else {
						o.WriteBody(resp.Body())
					}
				}
				return nil
			}())
		},
	}
	cmd.Flags().StringVar(&o.Description, "desc", "", "Description")
	return cmd
}
