package get

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

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

	// mcis
	cmds.AddCommand(&cobra.Command{
		Use:                "mcis (NAME | --name NAME) [options]",
		Short:              "Get MCISs.",
		Args:               app.BindCommandArgs(&o.Name),
		DisableFlagParsing: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, fnValidate())
			app.ValidateError(c, func() error {
				url := fmt.Sprintf("%s/ns/%s/mcis", app.Config.GetCurrentContext().Urls.Tumblebug, o.Namespace)
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

	// vmspec
	var config string
	cmdSpec := &cobra.Command{
		Use:   "spec --connection [name of connection-info.]",
		Short: "Get VM specifications.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if config == "" {
					c.Help()
				} else {
					// resty-go Get 인 경우 body를 혀용하지 않아서 "net/http" 모듈 사용
					body := bytes.NewBufferString(fmt.Sprintf("{\"connectionName\": \"%s\"}", config))
					req, err := http.NewRequest("GET", fmt.Sprintf("%s/vmspec", app.Config.GetCurrentContext().Urls.Spider), body)
					if err != nil {
						return err
					}
					req.Header.Add("Content-Type", "application/json")
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						return err
					}
					defer resp.Body.Close()

					bytes, _ := ioutil.ReadAll(resp.Body)
					o.WriteBody(bytes)
				}

				return nil
			}())
		},
	}
	cmdSpec.Flags().StringVar(&config, "connection", "", "Name of connection info.")
	cmds.AddCommand(cmdSpec)
	return cmds
}
