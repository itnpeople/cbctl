package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/itnpeople/cbctl/app"
)

// a struct to support command
type ConfigOptions struct {
	*app.Options
	Url_mcks      string
	Url_tumbelbug string
	Url_spider    string
}

func (o *ConfigOptions) writeYaml(in interface{}) {
	if b, err := yaml.Marshal(in); err != nil {
		o.PrintlnError(err)
	} else {
		o.WriteBody(b)
	}
}

// returns a cobra command
func NewCommandConfig(options *app.Options) *cobra.Command {
	o := &ConfigOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:   "config",
		Short: "configuration management",
		Long:  "",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// add-context
	cmdC := &cobra.Command{
		Use:                   "add-context [name]",
		Short:                 "Add a context",
		DisableFlagsInUseLine: true,
		Args:                  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if len(o.Name) == 0 {
					return fmt.Errorf("Name is required.")
				}
				if !strings.HasPrefix(o.Url_mcks, "http://") && !strings.HasPrefix(o.Url_mcks, "https://") {
					return fmt.Errorf("Invalid MCKS endpoint URL (value=%s)", o.Url_mcks)
				}
				if _, ok := app.Config.Contexts[o.Name]; ok {
					return fmt.Errorf("The context '%s' is alreaday exist", o.Name)
				} else {
					app.Config.Contexts[o.Name] = &app.ConfigContext{
						Name:      o.Name,
						Namespace: o.Namespace,
						Urls: struct {
							MCKS      string "yaml:\"mcks\""
							Spider    string "yaml:\"spider\""
							Tumblebug string "yaml:\"tumblebug\""
						}{MCKS: o.Url_mcks, Spider: o.Url_spider, Tumblebug: o.Url_tumbelbug},
					}
				}
				app.Config.WriteConfig()
				o.writeYaml(app.Config)
				return nil
			}())
		},
	}
	cmdC.Flags().StringVarP(&o.Url_mcks, "mcks", "", "", "MCKS endpoint URL (http://localhost:1470/mcks)")
	cmdC.Flags().StringVarP(&o.Url_tumbelbug, "tumblebug", "", "http://localhost:1323/tumblebug", "Tumblebug endpoint URL")
	cmdC.Flags().StringVarP(&o.Url_spider, "spider", "", "http://localhost:1024/spider", "Spider endpoint URL")
	cmds.AddCommand(cmdC)

	// view
	cmds.AddCommand(&cobra.Command{
		Use:   "view",
		Short: "Get contexts",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				o.writeYaml(app.Config)
				return nil
			}())
		},
	})

	// get context
	cmds.AddCommand(&cobra.Command{
		Use:   "get-context",
		Short: "Get a context",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					for k := range app.Config.Contexts {
						o.Println(k)
					}
				} else {
					if app.Config.Contexts[o.Name] != nil {
						o.writeYaml(app.Config.Contexts[o.Name])
					}
				}
				return nil
			}())
		},
	})

	// set context
	cmdS := &cobra.Command{
		Use:                   "set-context [name]",
		Short:                 "Set a context",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					c.Help()
				} else if app.Config.Contexts[o.Name] != nil {
					app.Config.Contexts[o.Name].Name = o.Name
					if o.Namespace != "" {
						app.Config.Contexts[o.Name].Namespace = o.Namespace
					}
					if o.Url_mcks != "" {
						app.Config.Contexts[o.Name].Urls.MCKS = o.Url_mcks
					}
					if o.Url_tumbelbug != "" {
						app.Config.Contexts[o.Name].Urls.Tumblebug = o.Url_tumbelbug
					}
					if o.Url_spider != "" {
						app.Config.Contexts[o.Name].Urls.Spider = o.Url_spider
					}
					o.writeYaml(app.Config.Contexts[o.Name])
				} else {
					o.Println("Not found a context (name=%s)", o.Name)
				}
				return nil
			}())
		},
	}
	cmdS.Flags().StringVarP(&o.Url_mcks, "mcks", "", "", "MCKS endpoint URL (http://localhost:1470/mcks)")
	cmdS.Flags().StringVarP(&o.Url_tumbelbug, "tumblebug", "", "http://localhost:1323/tumblebug", "Tumblebug endpoint URL")
	cmdS.Flags().StringVarP(&o.Url_spider, "spider", "", "http://localhost:1024/spider", "Spider endpoint URL")
	cmds.AddCommand(cmdS)

	// current-context (get/set)
	cmds.AddCommand(&cobra.Command{
		Use:                   "current-context [context]",
		Short:                 "Get/Set a current context",
		DisableFlagsInUseLine: true,
		Args:                  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if len(o.Name) > 0 {
					_, ok := app.Config.Contexts[o.Name]
					if ok {
						app.Config.CurrentContext = o.Name
						app.Config.WriteConfig()
					} else {
						o.Println("context '%s' is not exist\n", o.Name)
					}
				}
				o.writeYaml(app.Config.GetCurrentContext().Name)
				return nil
			}())
		},
	})

	// set-namespace
	cmds.AddCommand(&cobra.Command{
		Use:                   "set-namespace [namespace]",
		Short:                 "Get a current context",
		DisableFlagsInUseLine: true,
		Args:                  app.BindCommandArgs(&o.Namespace),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if len(o.Namespace) == 0 {
					c.Help()
				} else {
					app.Config.GetCurrentContext().Namespace = args[0]
					app.Config.WriteConfig()
					o.writeYaml(app.Config.GetCurrentContext())
				}
				return nil
			}())
		},
	})

	// delete-context
	cmds.AddCommand(&cobra.Command{
		Use:   "delete-context",
		Short: "Delete a context",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					return fmt.Errorf("Name Required.")
				}
				conf := app.Config
				if len(conf.Contexts) > 1 {
					delete(conf.Contexts, o.Name)
					if o.Name == conf.CurrentContext {
						conf.CurrentContext = func() string {
							if len(conf.Contexts) > 0 {
								for k := range conf.Contexts {
									return k
								}
							}
							return ""
						}()
					}
					conf.WriteConfig()
				}
				o.writeYaml(conf)
				return nil
			}())
		},
	})

	return cmds
}
