package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
)

// a struct to support command
type PluginOptions struct {
	app.IOStreams
}

// returns initialized Options
func newPluginOptions(ioStreams app.IOStreams) *PluginOptions {
	return &PluginOptions{
		IOStreams: ioStreams,
	}
}

// returns a cobra command
func NewCmdPlugin(streams app.IOStreams) *cobra.Command {
	o := NewVersionOptions(streams)
	cmd := &cobra.Command{
		Use:                   "plugin [flags]",
		DisableFlagsInUseLine: true,
		Short:                 "Provides utilities for interacting with plugins",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	cmd.AddCommand(NewCmdPluginList(o.IOStreams))
	return cmd
}

// a struct to support command
type PluginListOptions struct {
	app.IOStreams
	PluginPaths []string
	NameOnly    bool
	Verifier    PathVerifier
}

// returns initialized Options
func NewPluginListOptions(ioStreams app.IOStreams) *PluginListOptions {
	return &PluginListOptions{
		IOStreams: ioStreams,
	}
}

// returns a cobra command
func NewCmdPluginList(ioStreams app.IOStreams) *cobra.Command {

	o := NewPluginListOptions(ioStreams)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all visible plugin executables on a user's PATH",
		Run: func(cmd *cobra.Command, args []string) {
			app.ValidateError(o.Complete(cmd))
			app.ValidateError(o.Run())
		},
	}

	cmd.Flags().BoolVar(&o.NameOnly, "name-only", o.NameOnly, "If true, display only the binary name of each plugin, rather than its full path")
	return cmd
}

// completes all the required options
func (o *PluginListOptions) Complete(cmd *cobra.Command) error {
	o.Verifier = &CommandOverrideVerifier{
		root:        cmd.Root(),
		seenPlugins: make(map[string]string),
	}
	o.PluginPaths = filepath.SplitList(os.Getenv("PATH"))
	o.PluginPaths = append(o.PluginPaths, fmt.Sprintf("./%s", pluginDirectory))
	return nil
}

func (o *PluginListOptions) Run() error {
	pluginsFound := false
	isFirstFile := true
	pluginErrors := []error{}
	pluginWarnings := 0

	for _, dir := range uniquePathsList(o.PluginPaths) {
		if len(strings.TrimSpace(dir)) == 0 {
			continue
		}

		files, err := ioutil.ReadDir(dir)

		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				continue
			}
			continue
		}

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			if dir != "./"+pluginDirectory && !strings.HasPrefix(f.Name(), pluginFilenamePrefix+"-") {
				continue
			}

			if isFirstFile {
				fmt.Fprintf(o.Out, "The following compatible plugins are available:\n\n")
				pluginsFound = true
				isFirstFile = false
			}

			pluginPath := f.Name()
			if !o.NameOnly {
				pluginPath = filepath.Join(dir, pluginPath)
			}

			fmt.Fprintf(o.Out, "%s\n", pluginPath)
			if errs := o.Verifier.Verify(filepath.Join(dir, f.Name())); len(errs) != 0 {
				for _, err := range errs {
					fmt.Fprintf(o.ErrOut, "  - %s\n", err)
					pluginWarnings++
				}
			}
		}
	}

	if !pluginsFound {
		pluginErrors = append(pluginErrors, fmt.Errorf("error: unable to find any %s plugins in your PATH", pluginFilenamePrefix))
	}

	if pluginWarnings > 0 {
		if pluginWarnings == 1 {
			pluginErrors = append(pluginErrors, fmt.Errorf("error: one plugin warning was found"))
		} else {
			pluginErrors = append(pluginErrors, fmt.Errorf("error: %v plugin warnings were found", pluginWarnings))
		}
	}
	if len(pluginErrors) > 0 {
		errs := bytes.NewBuffer(nil)
		for _, e := range pluginErrors {
			fmt.Fprintln(errs, e)
		}
		return fmt.Errorf("%s", errs.String())
	}

	return nil
}

// receives a path and determines if it is valid or not
type PathVerifier interface {
	Verify(path string) []error
}

type CommandOverrideVerifier struct {
	root        *cobra.Command
	seenPlugins map[string]string
}

// implements PathVerifier and determines if a given path
func (v *CommandOverrideVerifier) Verify(path string) []error {
	if v.root == nil {
		return []error{fmt.Errorf("unable to verify path with nil root")}
	}

	// extract the plugin binary name
	segs := strings.Split(path, "/")
	binName := segs[len(segs)-1]

	cmdPath := strings.Split(binName, "-")
	if len(cmdPath) > 1 {
		// the first argument is always "cbctl" for a plugin binary
		cmdPath = cmdPath[1:]
	}

	errors := []error{}

	if isExec, err := isExecutable(path); err == nil && !isExec {
		errors = append(errors, fmt.Errorf("warning: %s identified as a %s plugin, but it is not executable", pluginFilenamePrefix, path))
	} else if err != nil {
		errors = append(errors, fmt.Errorf("error: unable to identify %s as an executable file: %v", path, err))
	}

	if existingPath, ok := v.seenPlugins[binName]; ok {
		errors = append(errors, fmt.Errorf("warning: %s is overshadowed by a similarly named plugin: %s", path, existingPath))
	} else {
		v.seenPlugins[binName] = path
	}

	if cmd, _, err := v.root.Find(cmdPath); err == nil {
		errors = append(errors, fmt.Errorf("warning: %s overwrites existing command: %q", binName, cmd.CommandPath()))
	}

	return errors
}

func isExecutable(fullPath string) (bool, error) {
	info, err := os.Stat(fullPath)
	if err != nil {
		return false, err
	}

	if runtime.GOOS == "windows" {
		fileExt := strings.ToLower(filepath.Ext(fullPath))

		switch fileExt {
		case ".bat", ".cmd", ".com", ".exe", ".ps1":
			return true, nil
		}
		return false, nil
	}

	if m := info.Mode(); !m.IsDir() && m&0111 != 0 {
		return true, nil
	}

	return false, nil
}

// deduplicates a given slice of strings without sorting or otherwise altering its order in any way.
func uniquePathsList(paths []string) []string {
	seen := map[string]bool{}
	newPaths := []string{}
	for _, p := range paths {
		if seen[p] {
			continue
		}
		seen[p] = true
		newPaths = append(newPaths, p)
	}
	return newPaths
}
