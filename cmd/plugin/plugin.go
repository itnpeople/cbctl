package plugin

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

const (
	PluginFilenamePrefix = "cbctl"
	PluginDirectory      = "plugins"
)

// a struct to support command
type PluginOptions struct {
	*app.Options
	PluginPaths          []string
	NameOnly             bool
	Verifier             PathVerifier
	pluginFilenamePrefix string
	pluginDirectory      string
}

// returns a cobra command
func NewCommandPlugin(options *app.Options) *cobra.Command {

	o := &PluginOptions{
		Options:              options,
		pluginFilenamePrefix: PluginFilenamePrefix,
		pluginDirectory:      PluginDirectory,
	}

	o.PluginPaths = filepath.SplitList(os.Getenv("PATH"))
	o.PluginPaths = append(o.PluginPaths, fmt.Sprintf("./%s", o.pluginDirectory))

	return &cobra.Command{
		Use:   "plugin",
		Short: "List all visible plugin executables on a user's PATH",
		Run: func(cmd *cobra.Command, args []string) {
			o.Verifier = &CommandOverrideVerifier{
				root:        cmd.Root(),
				seenPlugins: make(map[string]string),
			}
			o.PluginPaths = filepath.SplitList(os.Getenv("PATH"))
			o.PluginPaths = append(o.PluginPaths, fmt.Sprintf("./%s", o.pluginDirectory))
			app.ValidateError(cmd, o.Run())
		},
	}
}

func (o *PluginOptions) Run() error {
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
			if dir != "./"+o.pluginDirectory && !strings.HasPrefix(f.Name(), o.pluginFilenamePrefix+"-") {
				continue
			}

			if isFirstFile {
				o.Println("The following compatible plugins are available")
				pluginsFound = true
				isFirstFile = false
			}

			pluginPath := f.Name()
			if !o.NameOnly {
				pluginPath = filepath.Join(dir, pluginPath)
			}

			o.Println(pluginPath)
			if errs := o.Verifier.Verify(filepath.Join(dir, f.Name())); len(errs) != 0 {
				for _, err := range errs {
					o.PrintlnError(err)
					pluginWarnings++
				}
			}
		}
	}

	if !pluginsFound {
		pluginErrors = append(pluginErrors, fmt.Errorf("error: unable to find any %s plugins in your PATH", o.pluginFilenamePrefix))
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
		errors = append(errors, fmt.Errorf("warning: identified as a %s plugin, but it is not executable", path))
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
