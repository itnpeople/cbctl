package mcks

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type ClusterOptions struct {
	app.ConfigContext
	app.Output
	RootUrl      string
	Namespace    string
	Name         string
	Filenames    []string
	ControlPlane struct {
		Connection string
		Count      int
		Spec       string
	}
	Worker struct {
		Connection string
		Count      int
		Spec       string
	}
}

// returns initialized Options
func NewClusterOptions(ctx app.ConfigContext, output app.Output) *ClusterOptions {
	return &ClusterOptions{
		ConfigContext: ctx,
		Output:        output,
	}
}

// completes all the required options
func (o *ClusterOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, o.ConfigContext.Urls.MCKS)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request roo-url flag (%s)", o.RootUrl)
	}
	o.Namespace = utils.NVL(o.Namespace, o.ConfigContext.Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Invalid namespace flag")
	}
	return nil
}

// validates the provided options
func (o *ClusterOptions) Validate() error {
	if len(o.Filenames) == 0 {
		if o.ControlPlane.Connection == "" || o.Worker.Connection == "" {
			return fmt.Errorf("connection is empty")
		}
		if o.ControlPlane.Spec == "" || o.Worker.Spec == "" {
			return fmt.Errorf("spec is empty")
		}
	}
	return nil
}

// returns a cobra command
func NewCmdCluster(ctx app.ConfigContext, output app.Output) *cobra.Command {
	o := NewClusterOptions(ctx, output)

	// root
	cmds := &cobra.Command{
		Use:     "cluster",
		Short:   "Cluster management",
		Long:    "",
		Example: `  cbctl cluster`,
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "mcks root url (http://localhost:1470/mcks)")
	cmds.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "cloud-barista namespace for cluster list")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "cb-cluster", "cluster name")

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a cluster",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				var out []byte
				var err error
				if len(o.Filenames) > 0 {
					out, err = utils.ToFilenameBytes(o.Filenames[0])
				} else {
					out, err = utils.ToTemplateBytes(tplCluster, o)
				}
				if err != nil {
					return err
				}
				url := fmt.Sprintf("%s/ns/%s/clusters", o.RootUrl, o.Namespace)
				if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	}
	cmdC.Flags().StringVar(&o.ControlPlane.Connection, "control-plane-connection", "", "connection name of control-plane nodes")
	cmdC.Flags().IntVar(&o.ControlPlane.Count, "control-plane-count", 1, "count of control-plane nodes")
	cmdC.Flags().StringVar(&o.ControlPlane.Spec, "control-plane-spec", "", "spec of control-plane nodes")
	cmdC.Flags().StringVar(&o.Worker.Connection, "worker-connection", "", "connection name of wroker nodes")
	cmdC.Flags().IntVar(&o.Worker.Count, "worker-count", 1, "count of wroker nodes")
	cmdC.Flags().StringVar(&o.Worker.Spec, "worker-spec", "", "spec of wroker nodes")
	cmdC.Flags().StringArrayVarP(&o.Filenames, "filenames", "f", nil, "yaml")
	cmds.AddCommand(cmdC)

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all clusters",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/ns/%s/clusters", o.RootUrl, o.Namespace)); err != nil {
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
		Short: "Get a cluster",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/ns/%s/clusters/%s", o.RootUrl, o.Namespace, o.Name)); err != nil {
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
		Short: "Delete a cluster",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				// execute
				url := fmt.Sprintf("%s/ns/%s/clusters/%s", o.RootUrl, o.Namespace, o.Name)
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(url); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// update-kubeconfig
	cmds.AddCommand(&cobra.Command{
		Use:   "update-kubeconfig",
		Short: "Update a kubeconfig",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				url := fmt.Sprintf("%s/ns/%s/clusters/%s", o.RootUrl, o.Namespace, o.Name)
				res := &struct {
					ClusterConfig string `json:"clusterConfig"`
				}{}
				if _, err := resty.New().SetDisableWarn(true).R().SetResult(res).Get(url); err != nil {
					return err
				} else {
					conf, err := clientcmd.Load([]byte(res.ClusterConfig))
					if err != nil {
						return err
					}
					configLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
					cfg, err := configLoadingRules.Load()
					if err == nil {
						context := conf.Contexts[conf.CurrentContext]
						cfg.Clusters[fmt.Sprintf("%s-cluster", o.Name)] = conf.Clusters[context.Cluster]
						cfg.AuthInfos[fmt.Sprintf("%s-user", o.Name)] = conf.AuthInfos[context.AuthInfo]
						context.Cluster = fmt.Sprintf("%s-cluster", o.Name)
						context.AuthInfo = fmt.Sprintf("%s-user", o.Name)
						cfg.Contexts[o.Name] = context
						cfg.CurrentContext = o.Name
						err := clientcmd.WriteToFile(*cfg, configLoadingRules.GetDefaultFilename())
						if err != nil {
							return err
						}
					} else {
						return fmt.Errorf("cannot load kubeconfig (cause=%v)", err)
					}
					o.Output.Stream.WriteString("Success...\n")
				}
				return nil
			}())
		},
	})

	return cmds
}

const (
	tplCluster = `{
   "name": "{{.Name}}",
   "label": "",
   "description": "",
   "controlPlane": [
      { "connection": "{{.ControlPlane.Connection}}", "count": {{.ControlPlane.Count}}, "spec": "{{.ControlPlane.Spec}}" }
   ],
   "worker": [
      { "connection": "{{.Worker.Connection}}", "count": {{.Worker.Count}}, "spec": "{{.Worker.Spec}}" }
    ],
    "config": {
        "kubernetes": {
            "networkCni": "canal",
            "podCidr": "10.244.0.0/16",
            "serviceCidr": "10.96.0.0/12",
            "serviceDnsDomain": "cluster.local"
        }
    }
}`
)
