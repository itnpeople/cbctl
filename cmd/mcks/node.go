package mcks

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type NodeOptions struct {
	app.Output
	RootUrl   string
	Namespace string
	Cluster   string
	Name      string
	Worker    struct {
		Connection string
		Count      int
		Spec       string
	}
	Filenames []string
}

// returns initialized Options
func NewNodeOptions(output app.Output) *NodeOptions {
	return &NodeOptions{
		Output: output,
	}
}

// completes all the required options
func (o *NodeOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, app.Config.GetCurrentContext().Urls.MCKS)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid MCKS request-url flag (%s)", o.RootUrl)
	}
	o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Invalid namespace flag")
	}
	if o.Cluster == "" {
		return fmt.Errorf("Invalid cluster flag")
	}
	return nil
}

// validates the provided options
func (o *NodeOptions) Validate() error {
	if len(o.Filenames) == 0 {
		if o.Worker.Connection == "" {
			return fmt.Errorf("connection is empty")
		}
		if o.Worker.Spec == "" {
			return fmt.Errorf("spec is empty")
		}
	}
	return nil
}

// returns a cobra command
func NewCmdNodes(output app.Output) *cobra.Command {
	o := NewNodeOptions(output)

	// root
	cmds := &cobra.Command{
		Use:   "node",
		Short: "Nodes management",
		Long:  "",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "mcks root url (http://localhost:1470/mcks)")
	cmds.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "cloud-barista namespace for cluster list")
	cmds.PersistentFlags().StringVar(&o.Cluster, "cluster", "", "cluster name")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "node name")

	// create
	cmdC := &cobra.Command{
		Use:   "add",
		Short: "Add nodes",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				var out []byte
				var err error
				if len(o.Filenames) > 0 {
					out, err = utils.ToFilenameBytes(o.Filenames[0])
				} else {
					out, err = utils.ToTemplateBytes(tplNodes, o)
				}
				if err != nil {
					return err
				}
				url := fmt.Sprintf("%s/ns/%s/clusters/%s/nodes", o.RootUrl, o.Namespace, o.Cluster)
				if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	}
	cmdC.Flags().StringVar(&o.Worker.Connection, "worker-connection", "", "connection name of wroker nodes")
	cmdC.Flags().IntVar(&o.Worker.Count, "worker-count", 1, "count of wroker nodes")
	cmdC.Flags().StringVar(&o.Worker.Spec, "worker-spec", "", "spec of wroker nodes")
	cmdC.Flags().StringArrayVarP(&o.Filenames, "filenames", "f", nil, "yaml")
	cmds.AddCommand(cmdC)

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all nodes",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/ns/%s/clusters/%s/nodes", o.RootUrl, o.Namespace, o.Cluster)); err != nil {
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
		Short: "Get a node",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/ns/%s/clusters/%s/nodes/%s", o.RootUrl, o.Namespace, o.Cluster, o.Name)); err != nil {
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
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/ns/%s/clusters/%s/nodes/%s", o.RootUrl, o.Namespace, o.Cluster, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// get-key
	cmds.AddCommand(&cobra.Command{
		Use:   "get-key",
		Short: "Get a private key",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				res := &struct {
					Credential string `json:"credential"`
				}{}
				if _, err := resty.New().SetDisableWarn(true).R().SetResult(res).Get(fmt.Sprintf("%s/ns/%s/clusters/%s/nodes/%s", o.RootUrl, o.Namespace, o.Cluster, o.Name)); err != nil {
					return err
				} else {
					o.Output.Stream.WriteString(res.Credential)
				}
				return nil
			}())
		},
	})

	return cmds
}

const (
	tplNodes = `{
   "worker": [
      { "connection": "{{.Worker.Connection}}", "count": {{.Worker.Count}}, "spec": "{{.Worker.Spec}}" }
    ]
}`
)
