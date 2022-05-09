package getkey

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// returns a cobra command
func NewCommandGetKey(o *app.Options) *cobra.Command {

	// update-kubeconfig
	var clusterName string
	cmd := &cobra.Command{
		Use:                   "get-key (NAME | --name NAME) --cluster CLUSTER_NAME [options]",
		Short:                 "Get a SSH private key",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				// vlidation
				o.Namespace = utils.NVL(o.Namespace, app.Config.GetCurrentContext().Namespace)
				if o.Namespace == "" {
					return fmt.Errorf("Namespace is required.")
				}
				if o.Name == "" {
					return fmt.Errorf("Name is required.")
				}
				if clusterName == "" {
					return fmt.Errorf("Cluster name is required.")
				}

				// execute
				res := &struct {
					Credential string `json:"credential"`
				}{}
				url := fmt.Sprintf("%s/ns/%s/clusters/%s/nodes/%s", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace, clusterName, o.Name)
				if _, err := resty.New().SetDisableWarn(true).R().SetResult(res).Get(url); err != nil {
					return err
				} else {
					o.OutStream.WriteString(res.Credential)
				}
				return nil
			}())
		},
	}
	cmd.Flags().StringVar(&clusterName, "cluster", "", "Name of cluster")

	return cmd

}
