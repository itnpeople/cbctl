package create

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
)

// a struct to support command
type RegionOptions struct {
	*app.Options
	CSP           string
	Name          string
	Region        string
	Zone          string
	Location      string
	ResourceGroup string
}

// returns a cobra command
func NewCommandRegion(options *app.Options) *cobra.Command {
	o := &RegionOptions{
		Options: options,
	}

	cmd := &cobra.Command{
		Use:   "region",
		Short: "Create a cloud region",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.CSP == "" {
					return fmt.Errorf("CSP is required.")
				}
				if o.CSP == "azure" && (o.Location == "" || o.ResourceGroup == "") {
					return fmt.Errorf("Invalid location/resource-group flag (csp=%s, key=%s, secret=%s)", o.CSP, o.Location, o.ResourceGroup)
				}
				if out, err := app.GetBody(o, `{
					"RegionName"       : "{{ .Name }}",
					"ProviderName"     : "{{ .CSP | ToUpper }}", 
					"KeyValueInfoList" : [
						{"Key" : "Region",        "Value" : "{{ .Region }}"},
						{"Key" : "Zone",          "Value" : "{{ .Zone }}"},
						{"Key" : "location",      "Value" : "{{ .Location }}"},
						{"Key" : "ResourceGroup", "Value" : "{{ .ResourceGroup }}"}
					]
				}`); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/region", app.Config.GetCurrentContext().Urls.Spider)
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
	cmd.Flags().StringVar(&o.Region, "region", "", "Region")
	cmd.Flags().StringVar(&o.Zone, "zone", "", "Zone")
	cmd.Flags().StringVar(&o.Location, "location", "", "Location (azure)")
	cmd.Flags().StringVar(&o.ResourceGroup, "resource-group", "", "Resource group (azure)")

	return cmd
}
