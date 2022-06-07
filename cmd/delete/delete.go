package delete

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// returns a cobra command
func NewCommandDelete(o *app.Options) *cobra.Command {

	fnValidate := func() error {
		o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
		if o.Namespace == "" {
			return fmt.Errorf("Namespace is required.")
		}
		if o.Name == "" {
			return fmt.Errorf("Name is required.")
		}
		return nil
	}

	// root
	cmds := &cobra.Command{
		Use:                   "delete",
		Short:                 "Delete a object",
		DisableFlagsInUseLine: false,
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// cluster
	cmds.AddCommand(&cobra.Command{
		Use:                   "cluster (NAME | --name NAME) [options]",
		Short:                 "Delete a cluster",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: false,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/clusters/%s", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace, o.Name)
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// node
	var clusterName string
	cmdNode := &cobra.Command{
		Use:                   "node (NAME | --name NAME) --cluster CLUSTER_NAME [options]",
		Short:                 "Get nodes",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				if clusterName == "" {
					return fmt.Errorf("Cluster name is required.")
				}
				url := fmt.Sprintf("%s/ns/%s/clusters/%s/nodes/%s", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace, clusterName, o.Name)
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())

		},
	}
	cmdNode.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")
	cmds.AddCommand(cmdNode)

	// driver
	var csp string
	cmdDrv := &cobra.Command{
		Use:                   "driver (NAME | --name NAME) [options]",
		Short:                 "Delete a cloud driver",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/driver", app.Config.GetCurrentContext().Urls.Spider)
				if o.Name != "" {
					url += "/" + o.Name
				} else if csp != "" {
					url = fmt.Sprintf("%s/%s-driver-v1.0", url, csp)
				} else {
					return fmt.Errorf("Name is required.")
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	}
	cmdDrv.Flags().StringVar(&csp, "csp", "", "Cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack, cloudit)")
	cmds.AddCommand(cmdDrv)

	// region
	cmds.AddCommand(&cobra.Command{
		Use:                   "region (NAME | --name NAME) [options]",
		Short:                 "Delete a cloud region",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					return fmt.Errorf("Name is required.")
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/region/%s", app.Config.GetCurrentContext().Urls.Spider, o.Name)); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// credential
	cmds.AddCommand(&cobra.Command{
		Use:                   "credential (NAME | --name NAME) [options]",
		Short:                 "Delete a cloud credential",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					return fmt.Errorf("Name is required.")
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/credential/%s", app.Config.GetCurrentContext().Urls.Spider, o.Name)); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// connection
	cmds.AddCommand(&cobra.Command{
		Use:                   "connection (NAME | --name NAME) [options]",
		Short:                 "Delete a cloud connection info.",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					return fmt.Errorf("Name is required.")
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/connectionconfig/%s", app.Config.GetCurrentContext().Urls.Spider, o.Name)); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// namespace
	cmds.AddCommand(&cobra.Command{
		Use:                   "namespace (NAME | --name NAME) [options]",
		Short:                 "Delete a namespace.",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					return fmt.Errorf("Name is required.")
				}
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Delete(fmt.Sprintf("%s/ns/%s", app.Config.GetCurrentContext().Urls.Tumblebug, o.Name)); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// mcis
	cmds.AddCommand(&cobra.Command{
		Use:                   "mcis (NAME | --name NAME) [options]",
		Short:                 "Delete a MCIS.",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/mcis/%s?action=terminate", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace, o.Name)
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Delete(url + "?action=terminate"); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				if resp, err := http.Delete(url + "?action=refine"); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				if resp, err := http.Delete(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	return cmds
}
