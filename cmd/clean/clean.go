package clean

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type CleanOptions struct {
	*app.Options
}

// validates
func (o *CleanOptions) Validate() error {
	o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}
	return nil
}

func (o *CleanOptions) RunCleanupMCIS() error {

	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}

	url := fmt.Sprintf("%s/ns/%s", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
	http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")

	// mcis
	if resp, err := http.Delete(url + "/mcis"); err != nil {
		return err
	} else {
		o.WriteBody(resp.Body())
	}
	// vpc
	if resp, err := http.Delete(url + "/resources/vNet"); err != nil {
		return err
	} else {
		o.WriteBody(resp.Body())
	}
	// securityGroup
	if resp, err := http.Delete(url + "/resources/securityGroup"); err != nil {
		return err
	} else {
		o.WriteBody(resp.Body())
	}
	// sshKey
	if resp, err := http.Delete(url + "/resources/sshKey"); err != nil {
		return err
	} else {
		o.WriteBody(resp.Body())
	}
	// image
	if resp, err := http.Delete(url + "/resources/image"); err != nil {
		return err
	} else {
		o.WriteBody(resp.Body())
	}
	// spec
	if resp, err := http.Delete(url + "/resources/spec"); err != nil {
		return err
	} else {
		o.WriteBody(resp.Body())
	}
	return nil
}

// returns a cobra command
func NewCommandClean(options *app.Options) *cobra.Command {

	o := &CleanOptions{
		Options: options,
	}

	// clean
	cmd := &cobra.Command{
		Use:                   "clean [mcir]",
		Short:                 "Clean up objects",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, o.Validate())
			app.ValidateError(c, func() error {
				switch o.Name {
				case "mcir":
					return o.RunCleanupMCIS()
				default:
					c.Help()
					return nil
				}
			}())
		},
	}

	return cmd
}
