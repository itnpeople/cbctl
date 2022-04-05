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
type RegionOptions struct {
	app.Output
	RootUrl       string
	CSP           string
	Name          string
	Region        string
	Zone          string
	Location      string
	ResourceGroup string
}

// returns initialized Options
func NewRegionOptions(output app.Output) *RegionOptions {
	return &RegionOptions{
		Output: output,
	}
}

// completes all the required options
func (o *RegionOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, app.Config.GetCurrentContext().Urls.Spider)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request root-url flag (%s)", o.RootUrl)
	}
	return nil
}

// validates the provided options
func (o *RegionOptions) Validate() error {
	if o.CSP == "" {
		return fmt.Errorf("Invalid csp flag")
	}
	if o.CSP == "azure" && (o.Location == "" || o.ResourceGroup == "") {
		return fmt.Errorf("Invalid location/resource-group flag (csp=%s, key=%s, secret=%s)", o.CSP, o.Location, o.ResourceGroup)
	}
	return nil
}

// returns a cobra command
func NewCmdRegion(output app.Output) *cobra.Command {
	o := NewRegionOptions(output)
	cmds := &cobra.Command{
		Use:   "region",
		Short: "Cloud region",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "spider root url (http://localhost:1024/spider)")
	cmds.PersistentFlags().StringVar(&o.CSP, "csp", "", "cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack, cloudit)")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of region")

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a cloud region",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if out, err := utils.ToTemplateBytes(tplRegion, o); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/region", o.RootUrl)
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
	cmdC.Flags().StringVar(&o.Zone, "zone", "", "zone")
	cmdC.Flags().StringVar(&o.Location, "location", "", "location")
	cmdC.Flags().StringVar(&o.ResourceGroup, "resource-group", "", "resource group")
	cmds.AddCommand(cmdC)

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get cloud regions",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/region", o.RootUrl)); err != nil {
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
		Short: "Get a cloud region",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/region/%s", o.RootUrl, o.Name)); err != nil {
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
		Short: "Delete a cloud region",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/region/%s", o.RootUrl, o.Name)); err != nil {
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
	tplRegion = `{
"RegionName"       : "{{ .Name }}",
"ProviderName"     : "{{ .CSP | ToUpper }}", 
"KeyValueInfoList" : [
	{"Key" : "Region",        "Value" : "{{ .Region }}"},
	{"Key" : "Zone",          "Value" : "{{ .Zone }}"},
	{"Key" : "location",      "Value" : "{{ .Location }}"},
	{"Key" : "ResourceGroup", "Value" : "{{ .ResourceGroup }}"}
]}`
)
