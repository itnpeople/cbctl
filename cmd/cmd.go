package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/cmd/clean"
	"github.com/itnpeople/cbctl/cmd/config"
	"github.com/itnpeople/cbctl/cmd/create"
	"github.com/itnpeople/cbctl/cmd/delete"
	"github.com/itnpeople/cbctl/cmd/get"
	"github.com/itnpeople/cbctl/cmd/get-key"
	"github.com/itnpeople/cbctl/cmd/plugin"
	"github.com/itnpeople/cbctl/cmd/update-kubeconfig"
)

func Execute() {
	rootCmd := NewRootCommand()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type CBCtlOptions struct {
	app.Options
	plugin.PluginHandler
}

func NewRootCommand() *cobra.Command {

	o := CBCtlOptions{
		PluginHandler: plugin.NewDefaultPluginHandler(),
		Options: app.Options{
			OutStream: os.Stdout,
		},
	}

	// cbctl
	cmds := &cobra.Command{
		Use:                   "cbctl",
		Short:                 "Cloud-Barista. command-line-interface manager",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	// Persistent Flags
	cmds.PersistentFlags().StringVarP(&o.ConfigFile, "config", "c", "", "Configuration file path")
	cmds.PersistentFlags().StringVarP(&o.Output, "output", "o", app.OUTPUT_YAML, "Output format(yaml/json)")
	cmds.PersistentFlags().StringVarP(&o.Filename, "file", "f", "", "Filename")
	cmds.PersistentFlags().StringVarP(&o.Namespace, "namespace", "n", "", "Cloud-barista namespace")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "Name")

	// initialize config file
	if err := app.OnConfigInitialize(o.ConfigFile); err != nil {
		o.PrintlnError(err)
		os.Exit(1)
	}

	// add commands
	cmds.AddCommand(&cobra.Command{
		Use:                   "version",
		Short:                 "Print the version number of cbctl",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			o.Println("Version=%s, buildTime=%s", app.BuildVersion, app.BuildTime)
		},
	})
	cmds.AddCommand(get.NewCommandGet(&o.Options))                           // cbctl get
	cmds.AddCommand(create.NewCommandCreate(&o.Options))                     // cbctl create
	cmds.AddCommand(delete.NewCommandDelete(&o.Options))                     // cbctl delete
	cmds.AddCommand(config.NewCommandConfig(&o.Options))                     // cbctl config
	cmds.AddCommand(updatekubeconfig.NewCommandUpdateKubeconfig(&o.Options)) // cbctl update-kubeconfig
	cmds.AddCommand(getkey.NewCommandGetKey(&o.Options))                     // cbctl get-key
	cmds.AddCommand(plugin.NewCommandPlugin(&o.Options))                     // cbctl plugin
	cmds.AddCommand(clean.NewCommandClean(&o.Options))                       // cbctl clean

	// execute plugin
	if len(os.Args) > 0 {
		cmdPathPieces := os.Args[1:]
		if _, _, err := cmds.Find(cmdPathPieces); err != nil {
			if err := plugin.HandlePluginCommand(o.PluginHandler, cmdPathPieces); err != nil {
				o.PrintlnError(err)
				os.Exit(1)
			}
		}
	}

	return cmds

}
