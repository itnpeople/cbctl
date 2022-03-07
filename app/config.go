package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	CurrentContext string                    `yaml:"current-context"`
	Contexts       map[string]*ConfigContext `yaml:"contexts"`
}

type ConfigContext struct {
	Namespace string `yaml:"namespace"`
	Urls      struct {
		MCKS      string `yaml:"mcks"`
		Spider    string `yaml:"spider"`
		Tumblebug string `yaml:"tumblebug"`
	} `yaml:"urls"`
}

func GetConfig(cfgFile *string) (*Config, error) {
	dir := fmt.Sprintf("%s/.cbctl", HomeDir())
	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if *cfgFile != "" {
		viper.SetConfigFile(*cfgFile)
	}

	// read a config file
	if err := viper.ReadInConfig(); err != nil {

		// the default config save to "${HOME}/.cbctl/config"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
		if _, err := os.Stat(filepath.Join(dir, "config")); os.IsNotExist(err) {
			ioutil.WriteFile(filepath.Join(dir, "config"), []byte(`current-context: local
contexts:
  local:
    namespace :
    urls:
      mcks: http://localhost:1470/mcks
      spider: http://localhost:1024/spider
      tumblebug: http://localhost:1323/tumblebug`), os.ModePerm)
		}

		if err = viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("Fail to using config file: %s (cause=%v)", viper.ConfigFileUsed(), err)
		}
	}

	conf := &Config{}
	if err := viper.Unmarshal(&conf,
		viper.DecoderConfigOption(func(decoderConfig *mapstructure.DecoderConfig) {
			decoderConfig.TagName = "yaml"
		})); err != nil {
		return nil, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	return conf, nil

}

func GetCurrentContext(cfgFile *string) (*ConfigContext, error) {

	if conf, err := GetConfig(cfgFile); err != nil {
		return nil, err
	} else {
		if conf.Contexts[conf.CurrentContext] == nil {
			conf.CurrentContext = func() string {
				if len(conf.Contexts) > 0 {
					for k := range conf.Contexts {
						return k
					}
				}
				return ""
			}()
		}
		if conf.CurrentContext == "" {
			return nil, fmt.Errorf("unable to find current context")
		}

		return conf.Contexts[conf.CurrentContext], nil
	}
}

func HomeDir() string {

	if runtime.GOOS == "windows" {
		home := os.Getenv("HOME")
		homeDriveHomePath := ""
		if homeDrive, homePath := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"); len(homeDrive) > 0 && len(homePath) > 0 {
			homeDriveHomePath = homeDrive + homePath
		}
		userProfile := os.Getenv("USERPROFILE")

		// Return first of %HOME%, %HOMEDRIVE%/%HOMEPATH%, %USERPROFILE% that contains a `.cbctl\config` file.
		// %HOMEDRIVE%/%HOMEPATH% is preferred over %USERPROFILE% for backwards-compatibility.
		for _, p := range []string{home, homeDriveHomePath, userProfile} {
			if len(p) == 0 {
				continue
			}
			if _, err := os.Stat(filepath.Join(p, ".cbctl", "config")); err != nil {
				continue
			}
			return p
		}

		firstSetPath := ""
		firstExistingPath := ""

		// Prefer %USERPROFILE% over %HOMEDRIVE%/%HOMEPATH% for compatibility with other auth-writing tools
		for _, p := range []string{home, userProfile, homeDriveHomePath} {
			if len(p) == 0 {
				continue
			}
			if len(firstSetPath) == 0 {
				// remember the first path that is set
				firstSetPath = p
			}
			info, err := os.Stat(p)
			if err != nil {
				continue
			}
			if len(firstExistingPath) == 0 {
				// remember the first path that exists
				firstExistingPath = p
			}
			if info.IsDir() && info.Mode().Perm()&(1<<(uint(7))) != 0 {
				// return first path that is writeable
				return p
			}
		}

		// If none are writeable, return first location that exists
		if len(firstExistingPath) > 0 {
			return firstExistingPath
		}

		// If none exist, return first location that is set
		if len(firstSetPath) > 0 {
			return firstSetPath
		}

		// We've got nothing
		return ""
	}
	return os.Getenv("HOME")
}
