package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/itnpeople/cbctl/app"
)

type PluginHandler interface {
	Lookup(filename string) (string, bool)
	Execute(executablePath string, cmdArgs, environment []string) error
}

type DefaultPluginHandler struct {
	prefix             string
	pluginsDirectories []string
}

func NewDefaultPluginHandler() *DefaultPluginHandler {

	return &DefaultPluginHandler{
		prefix:             PluginFilenamePrefix,
		pluginsDirectories: []string{filepath.Join(app.HomeDir(), ".cbctl", PluginDirectory)},
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
