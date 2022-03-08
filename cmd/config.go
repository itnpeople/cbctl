package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type ConfigOptions struct {
	app.Output
	Name          string
	Namespace     string
	Url_mcks      string
	Url_tumbelbug string
	Url_spider    string
}

// returns initialized Options
func NewConfigOptions(output app.Output) *ConfigOptions {
	return &ConfigOptions{
		Output: output,
	}
}

// completes all the required options
func (o *ConfigOptions) Complete(cmd *cobra.Command) error {
	return nil
}

// validates the provided options
func (o *ConfigOptions) Validate() error {
	if len(o.Name) == 0 {
		return fmt.Errorf("Invalid paremeters (name=%s)", o.Name)
	}
	if !strings.HasPrefix(o.Url_mcks, "http://") && !strings.HasPrefix(o.Url_mcks, "https://") {
		return fmt.Errorf("Invalid paremeters (url-mcks=%s)", o.Url_mcks)
	}
	if !strings.HasPrefix(o.Url_spider, "http://") && !strings.HasPrefix(o.Url_spider, "https://") {
		return fmt.Errorf("Invalid paremeters (url-spider=%s)", o.Url_spider)
	}
	return nil
}

// returns a cobra command
func NewCmdConfig(output app.Output) *cobra.Command {
	o := NewConfigOptions(output)

	// root
	cmds := &cobra.Command{
		Use:   "config",
		Short: "configuration management",
		Long:  "",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "context name")

	// add-context
	cmdC := &cobra.Command{
		Use:   "add-context",
		Short: "Add a context",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				fmt.Println(app.Config)
				if _, ok := app.Config.Contexts[o.Name]; ok {
					return fmt.Errorf("the context '%s' is alreaday exist", o.Name)
				} else {
					app.Config.Contexts[o.Name] = &app.ConfigContext{
						Namespace: o.Namespace,
						Urls: struct {
							MCKS      string "yaml:\"mcks\""
							Spider    string "yaml:\"spider\""
							Tumblebug string "yaml:\"tumblebug\""
						}{MCKS: o.Url_mcks, Spider: o.Url_spider, Tumblebug: o.Url_tumbelbug},
					}
				}
				app.Config.WriteConfig()
				outYaml(app.Config, o.Output)
				return nil
			}())
		},
	}
	cmdC.Flags().StringVar(&o.Namespace, "namespace", "", "context default namespace")
	cmdC.Flags().StringVar(&o.Url_mcks, "url-mcks", "http://localhost:1470/mcks", "mcks rest-root-url (http://localhost:1470/mcks)")
	cmdC.Flags().StringVar(&o.Url_tumbelbug, "url-tumblebug", "http://localhost:1323/tumblebug", "tumblebug rest-root-url (http://localhost:1323/tumblebug)")
	cmdC.Flags().StringVar(&o.Url_spider, "url-spider", "http://localhost:1024/spider", "spider  rest-root-url (http://localhost:1024/spider)")
	cmds.AddCommand(cmdC)

	// view
	cmds.AddCommand(&cobra.Command{
		Use:   "view",
		Short: "Get contexts",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				outYaml(app.Config, o.Output)
				return nil
			}())
		},
	})

	// context-list
	cmds.AddCommand(&cobra.Command{
		Use:   "list-context",
		Short: "Get contexts",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				for k := range app.Config.Contexts {
					o.Output.WriteString(k + "\n")
				}
				return nil
			}())
		},
	})

	// context-get
	cmds.AddCommand(&cobra.Command{
		Use:   "get-context",
		Short: "List all clusters",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if app.Config.Contexts[o.Name] != nil {
					outYaml(app.Config.Contexts[o.Name], o.Output)
				}
				return nil
			}())
		},
	})

	// current-context (get/set)
	cmds.AddCommand(&cobra.Command{
		Use:   "current-context",
		Short: "Get a current context",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
					_, ok := app.Config.Contexts[o.Name]
					if ok {
						app.Config.CurrentContext = o.Name
						app.Config.WriteConfig()
					} else {
						o.Output.WriteString("context '%s' is not exist\n", o.Name)
					}
				} else {
					o.Output.WriteString(app.Config.CurrentContext + "\n")
				}
				return nil
			}())
		},
	})

	// set-namespace
	cmds.AddCommand(&cobra.Command{
		Use:   "set-namespace",
		Short: "Get a current context",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					app.Config.GetCurrentContext().Namespace = args[0]
					app.Config.WriteConfig()
				}
				outYaml(app.Config.GetCurrentContext(), o.Output)
				return nil
			}())
		},
	})

	// delete-context
	cmds.AddCommand(&cobra.Command{
		Use:   "delete-context",
		Short: "Delete a context",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
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
				outYaml(conf, o.Output)
				return nil
			}())
		},
	})

	return cmds
}

func outYaml(in interface{}, out app.Output) {
	if b, err := yaml.Marshal(in); err != nil {
		out.Stream.WriteString(err.Error())
	} else {
		out.Write(b)
	}
}
