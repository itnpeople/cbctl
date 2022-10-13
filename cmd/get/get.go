package get

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// returns a cobra command
func NewCommandGet(o *app.Options) *cobra.Command {

	// validate namespace
	fnValidate := func() error {
		o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
		if o.Namespace == "" {
			return fmt.Errorf("Namespace is required.")
		}
		return nil
	}

	// Get command
	cmds := &cobra.Command{
		Use:                   "get",
		Short:                 "Get objects",
		DisableFlagsInUseLine: false,
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// get cluster command
	cmds.AddCommand(&cobra.Command{
		Use:                   "cluster (NAME | --name NAME) [options]",
		Short:                 "Get clusters",
		DisableFlagsInUseLine: false,
		Args:                  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/clusters", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace)
				if o.Name != "" {
					url = url + "/" + o.Name
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// get nodes command
	var clusterName string
	cmdNode := &cobra.Command{
		Use:                   "node (NAME | --name NAME) [options]",
		Short:                 "Get nodes",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				if clusterName == "" {
					return fmt.Errorf("Cluster name is required.")
				}
				url := fmt.Sprintf("%s/ns/%s/clusters/%s/nodes", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace, clusterName)
				if o.Name != "" {
					url = url + "/" + o.Name
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(url); err != nil {
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
		Short:                 "Get cloud drivers",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/driver", app.Config.GetCurrentContext().Urls.Spider)
				if o.Name != "" {
					url = fmt.Sprintf("%s/%s", url, o.Name)
				} else if csp != "" {
					url = fmt.Sprintf("%s/%s-driver-v1.0", url, csp)
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(url); err != nil {
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
		Short:                 "Get cloud regions",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/region", app.Config.GetCurrentContext().Urls.Spider)
				if o.Name != "" {
					url += "/" + o.Name
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(url); err != nil {
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
		Short:                 "Get a cloud credential",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/credential", app.Config.GetCurrentContext().Urls.Spider)
				if o.Name != "" {
					url += "/" + o.Name
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// connection info.
	cmds.AddCommand(&cobra.Command{
		Use:                   "connection (NAME | --name NAME) [options]",
		Short:                 "Get a cloud connection infos.",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/connectionconfig", app.Config.GetCurrentContext().Urls.Spider)
				if o.Name != "" {
					url += "/" + o.Name
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(url); err != nil {
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
		Short:                 "Get cloud-barista namespaces.",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns", app.Config.GetCurrentContext().Urls.Tumblebug)
				if o.Name != "" {
					url += "/" + o.Name
				}
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// vpc
	cmds.AddCommand(&cobra.Command{
		Use:   "vpc (NAME | --name NAME) [options]",
		Short: "Get VPCs.",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/resources/vNet", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
				if o.Name != "" {
					url += "/" + o.Name
				}
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}

				return nil
			}())
		},
	})

	// security group
	cmds.AddCommand(&cobra.Command{
		Use:   "sg (NAME | --name NAME) [options]",
		Short: "Get Security Groups.",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/resources/securityGroup", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
				if o.Name != "" {
					url += "/" + o.Name
				}
				fmt.Println(url)
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// ssh-key
	cmds.AddCommand(&cobra.Command{
		Use:   "sshkey (NAME | --name NAME) [options]",
		Short: "Get SSH Keys.",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/resources/sshKey", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
				if o.Name != "" {
					url += "/" + o.Name
				}
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// images
	cmds.AddCommand(&cobra.Command{
		Use:   "image (NAME | --name NAME) [options]",
		Short: "Get Disk Images.",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/resources/image", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
				if o.Name != "" {
					url += "/" + o.Name
				}
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
					return err
				} else {
					o.WriteBody(resp.Body())
				}
				return nil
			}())
		},
	})

	// spec
	cmds.AddCommand(&cobra.Command{
		Use:   "spec (NAME | --name NAME) [options]",
		Short: "Get VM specifications.",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/resources/spec", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
				if o.Name != "" {
					url += "/" + o.Name
				}
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
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
		Use:   "mcis (NAME | --name NAME) [options]",
		Short: "Get MCISs.",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/mcis", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
				if o.Name != "" {
					url += "/" + o.Name
				}
				fmt.Println(url)
				http := resty.New().SetDisableWarn(true).R().SetBasicAuth("default", "default")
				if resp, err := http.Get(url); err != nil {
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
