package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	//"github.com/spf13/viper"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/cmd/mcks"
	"github.com/itnpeople/cbctl/cmd/spider"
)

const (
	pluginFilenamePrefix = "cbctl"
	pluginDirectory      = "plugins"
)

var (
	rootCmd *cobra.Command
)

func Execute() {
	rootCmd = NewDefaultCBCtlCommand()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type CBCtlOptions struct {
	app.ConfigContext
	app.IOStreams
	PluginHandler

	Arguments []string
	Output    string
}

func NewDefaultCBCtlCommand() *cobra.Command {
	return NewDefaultCBCtlCommandWithArgs(CBCtlOptions{
		PluginHandler: NewDefaultPluginHandler(pluginFilenamePrefix, []string{filepath.Join(app.HomeDir(), ".cbctl", pluginDirectory)}),
		Arguments:     os.Args,
		IOStreams:     app.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

func NewDefaultCBCtlCommandWithArgs(o CBCtlOptions) *cobra.Command {
	cmd := NewCBCtlCommand(o)

	if o.PluginHandler == nil {
		return cmd
	}

	if len(o.Arguments) > 1 {
		cmdPathPieces := o.Arguments[1:]
		if _, _, err := cmd.Find(cmdPathPieces); err != nil {
			var cmdName string
			for _, arg := range cmdPathPieces {
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}
			switch cmdName {
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			default:
				if err := HandlePluginCommand(o.PluginHandler, cmdPathPieces); err != nil {
					fmt.Fprintf(o.IOStreams.ErrOut, "Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	return cmd
}

func NewCBCtlCommand(o CBCtlOptions) *cobra.Command {

	cmds := &cobra.Command{
		Use:   "cbctl",
		Short: "Cloud-Barista. command-line-interface manager",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var cfgFile string
	cmds.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	cmds.PersistentFlags().StringVarP(&o.Output, "output", "o", app.OUTPUT_YAML, "output (yaml or json)")

	output := app.Output{Type: &o.Output, Stream: o.IOStreams.Out}
	cmds.AddCommand(NewCmdVersion(o.IOStreams))
	cmds.AddCommand(mcks.NewCmdCluster(o.ConfigContext, output))
	cmds.AddCommand(mcks.NewCmdNodes(output))
	cmds.AddCommand(spider.NewCmdDriver(o.ConfigContext, output))
	cmds.AddCommand(spider.NewCmdCredential(o.ConfigContext, output))
	cmds.AddCommand(spider.NewCmdRegion(o.ConfigContext, output))
	cmds.AddCommand(spider.NewCmdConnection(o.ConfigContext, output))
	cmds.AddCommand(NewCmdPlugin(o.IOStreams))

	cobra.OnInitialize(func() {
		if ctx, err := app.GetCurrentContext(&cfgFile); err != nil {
			fmt.Println(err.Error())
		} else {
			o.ConfigContext = *ctx
		}
	})

	return cmds
}

type PluginHandler interface {
	Lookup(filename string) (string, bool)
	Execute(executablePath string, cmdArgs, environment []string) error
}

type DefaultPluginHandler struct {
	prefix             string
	pluginsDirectories []string
}

func NewDefaultPluginHandler(prefix string, dirs []string) *DefaultPluginHandler {
	return &DefaultPluginHandler{
		prefix:             prefix,
		pluginsDirectories: dirs,
	}
}

func (h *DefaultPluginHandler) Lookup(filename string) (string, bool) {

	found := false
	path, err := exec.LookPath(fmt.Sprintf("%s-%s", h.prefix, filename))
	if err != nil || len(path) == 0 {
		for _, dir := range h.pluginsDirectories {
			path, err = exec.LookPath(filepath.Join(dir, filename))
			if err == nil && len(path) > 0 {
				found = true
				break
			}
		}
	} else {
		found = true
	}
	return path, found

}

func (h *DefaultPluginHandler) Execute(executablePath string, cmdArgs, environment []string) error {

	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(executablePath, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = environment
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}

	// invoke cmd binary relaying the environment and args given
	return syscall.Exec(executablePath, append([]string{executablePath}, cmdArgs...), environment)
}

func HandlePluginCommand(pluginHandler PluginHandler, cmdArgs []string) error {
	var remainingArgs []string // all "non-flag" arguments
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, "-") {
			break
		}
		remainingArgs = append(remainingArgs, strings.Replace(arg, "-", "_", -1))
	}

	if len(remainingArgs) == 0 {
		// the length of cmdArgs is at least 1
		return fmt.Errorf("flags cannot be placed before plugin name: %s\n", cmdArgs[0])
	}

	foundBinaryPath := ""
	for len(remainingArgs) > 0 {
		path, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
		if !found {
			remainingArgs = remainingArgs[:len(remainingArgs)-1]
			continue
		}

		foundBinaryPath = path
		break
	}

	if len(foundBinaryPath) == 0 {
		return nil
	}

	// invoke cmd binary relaying the current environment and args given
	if err := pluginHandler.Execute(foundBinaryPath, cmdArgs[len(remainingArgs):], os.Environ()); err != nil {
		return err
	}

	return nil
}
