package create

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type CreateClusterOptions struct {
	*app.Options
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

// validates
func (o *CreateClusterOptions) Validate() error {
	o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}
	return nil
}

func (o *CreateClusterOptions) Run() error {

	// execute
	if out, err := app.GetBody(o, `{
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
		}`); err != nil {
		return err
	} else {
		url := fmt.Sprintf("%s/ns/%s/clusters", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace)
		if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
			return err
		} else {
			o.WriteBody(resp.Body())
		}
	}
	return nil
}

// a struct to support command
type CreateNodeOptions struct {
	*app.Options
	clusterName string
	Worker      struct {
		Connection string
		Count      int
		Spec       string
	}
}

// validates
func (o *CreateNodeOptions) Validate() error {
	o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
	if o.Namespace == "" {
		return fmt.Errorf("Namespace is required.")
	}
	if o.clusterName == "" {
		return fmt.Errorf("Cluster is required.")
	}
	return nil
}

func (o *CreateNodeOptions) Run() error {

	// exute
	if out, err := app.GetBody(o, `
		{"worker": [
			{ "connection": "{{.Worker.Connection}}", "count": {{.Worker.Count}}, "spec": "{{.Worker.Spec}}" }
		]}`); err != nil {
		return err
	} else {
		url := fmt.Sprintf("%s/ns/%s/clusters/%s/nodes", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace, o.clusterName)
		if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
			return err
		} else {
			o.WriteBody(resp.Body())
		}
	}
	return nil

}

// returns a cobra command
func NewCommandCreate(options *app.Options) *cobra.Command {
	oCluster := &CreateClusterOptions{
		Options: options,
	}

	// cbctl create
	cmds := &cobra.Command{
		Use:   "create",
		Short: "Create a object",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// cbctl create cluster
	cmdC := &cobra.Command{
		Use:   "cluster",
		Short: "Create a cluster",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, oCluster.Validate())
			app.ValidateError(c, oCluster.Run())
		},
	}
	cmdC.Flags().StringVar(&oCluster.ControlPlane.Connection, "control-plane-connection", "", "Connection name of control-plane nodes")
	cmdC.Flags().IntVar(&oCluster.ControlPlane.Count, "control-plane-count", 1, "Count of control-plane nodes")
	cmdC.Flags().StringVar(&oCluster.ControlPlane.Spec, "control-plane-spec", "", "Spec. of control-plane nodes")
	cmdC.Flags().StringVar(&oCluster.Worker.Connection, "worker-connection", "", "Connection name of wroker nodes")
	cmdC.Flags().IntVar(&oCluster.Worker.Count, "worker-count", 1, "Count of wroker nodes")
	cmdC.Flags().StringVar(&oCluster.Worker.Spec, "worker-spec", "", "Spec. of wroker nodes")
	cmds.AddCommand(cmdC)

	oNode := &CreateNodeOptions{
		Options: options,
	}

	// cbctl create node
	cmdN := &cobra.Command{
		Use:   "node",
		Short: "Add nodes",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, oNode.Validate())
			app.ValidateError(c, oNode.Run())
		},
	}
	cmdN.Flags().StringVar(&oNode.clusterName, "cluster", "", "Name of cluster")
	cmdN.Flags().StringVar(&oNode.Worker.Connection, "worker-connection", "", "Connection name of wroker nodes")
	cmdN.Flags().IntVar(&oNode.Worker.Count, "worker-count", 1, "Count of wroker nodes")
	cmdN.Flags().StringVar(&oNode.Worker.Spec, "worker-spec", "", "Spec. of wroker nodes")
	cmds.AddCommand(cmdN)

	cmds.AddCommand(NewCommandDriver(options))     // cbctl crate driver
	cmds.AddCommand(NewCommandRegion(options))     // cbctl create region
	cmds.AddCommand(NewCommandCredential(options)) // cbctl create credential
	cmds.AddCommand(NewCommandConnection(options)) // cbctl create conenection
	cmds.AddCommand(NewCommandConnection(options)) // cbctl create conenection
	cmds.AddCommand(NewCommandNamespace(options))  // cbctl create namespace

	return cmds
}
