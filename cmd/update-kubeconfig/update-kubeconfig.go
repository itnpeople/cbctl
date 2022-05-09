package updatekubeconfig

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// options
type UpdateKubeconfigOptions struct {
	*app.Options
}

// returns a cobra command
func NewCommandUpdateKubeconfig(options *app.Options) *cobra.Command {
	o := &UpdateKubeconfigOptions{
		Options: options,
	}

	// update-kubeconfig
	return &cobra.Command{
		Use:                   "update-kubeconfig (NAME | --name NAME) [options]",
		Short:                 "Update a kubeconfig",
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

				// execute
				url := fmt.Sprintf("%s/ns/%s/clusters/%s", app.Config.GetCurrentContext().Urls.MCKS, o.Namespace, o.Name)
				res := &struct {
					ClusterConfig string `json:"clusterConfig"`
				}{}
				if resp, err := resty.New().SetDisableWarn(true).R().SetResult(res).Get(url); err != nil {
					return err
				} else if res.ClusterConfig != "" {
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
					o.Println("Success...")
				} else {
					o.WriteBody(resp.Body())
				}
				return nil

			}())
		},
	}

}
