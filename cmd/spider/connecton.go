package spider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"github.com/itnpeople/cbctl/app"
	"github.com/itnpeople/cbctl/utils"
)

// a struct to support command
type ConnectionOptions struct {
	app.Output
	RootUrl    string
	Name       string
	CSP        string
	Credential string
	Region     string
}

// returns initialized Options
func NewConnectionOptions(output app.Output) *ConnectionOptions {
	return &ConnectionOptions{
		Output: output,
	}
}

// completes all the required options
func (o *ConnectionOptions) Complete(cmd *cobra.Command) error {
	o.RootUrl = utils.NVL(o.RootUrl, app.Config.GetCurrentContext().Urls.Spider)
	if !strings.HasPrefix(o.RootUrl, "http://") && !strings.HasPrefix(o.RootUrl, "https://") {
		return fmt.Errorf("Invalid request root-url flag (%s)", o.RootUrl)
	}
	return nil
}

// validates the provided options
func (o *ConnectionOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("Invalid name flag")
	}
	if o.CSP == "" {
		return fmt.Errorf("Invalid csp flag")
	}
	if o.Credential == "" {
		return fmt.Errorf("Invalid credential flag")
	}
	if o.Region == "" {
		return fmt.Errorf("Invalid region flag")
	}
	return nil
}

// returns a cobra command
func NewCmdConnection(output app.Output) *cobra.Command {
	o := NewConnectionOptions(output)
	cmds := &cobra.Command{
		Use:   "connection",
		Short: "Cloud connection info.",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmds.PersistentFlags().StringVar(&o.RootUrl, "url", "", "spider root url (http://localhost:1024/spider)")
	cmds.PersistentFlags().StringVar(&o.CSP, "csp", "", "cloud service provider (aws, gcp, azure, alibaba, tencent, ibm, openstack, cloudit)")
	cmds.PersistentFlags().StringVar(&o.Name, "name", "", "name of connection info.")

	// create
	cmdC := &cobra.Command{
		Use:   "create",
		Short: "Create a cloud connection info.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(o.Validate())
			app.ValidateError(func() error {
				if out, err := utils.ToTemplateBytes(tplConnection, o); err != nil {
					return err
				} else {
					url := fmt.Sprintf("%s/connectionconfig", o.RootUrl)
					if resp, err := resty.New().SetDisableWarn(true).R().SetHeader("content-type", "application/json").SetBody(out).Post(url); err != nil {
						return err
					} else {
						o.Output.Write(resp.Body())
					}
				}
				return nil
			}())
		},
	}
	cmdC.Flags().StringVar(&o.Region, "region", "", "region")
	cmdC.Flags().StringVar(&o.Credential, "credential", "", "credential")
	cmds.AddCommand(cmdC)

	// list
	cmds.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Get cloud connection infos.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/connectionconfig", o.RootUrl)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// get
	cmds.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a cloud connection infos.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Get(fmt.Sprintf("%s/connectionconfig/%s", o.RootUrl, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// delete
	cmds.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete a cloud connection info.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if resp, err := resty.New().SetDisableWarn(true).R().Delete(fmt.Sprintf("%s/connectionconfig/%s", o.RootUrl, o.Name)); err != nil {
					return err
				} else {
					o.Output.Write(resp.Body())
				}
				return nil
			}())
		},
	})

	// test
	cmds.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "Test a cloud connection infos.",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(o.Complete(c))
			app.ValidateError(func() error {
				if len(args) > 0 {
					o.Name = utils.NVL(o.Name, args[0])
				}
				if o.Name == "" {
					return fmt.Errorf("Invalid name flag")
				}
				// resty-go Get 인 경우 body를 혀용하지 않아서 "net/http" 모듈 사용
				body := bytes.NewBufferString(fmt.Sprintf("{\"connectionName\": \"%s\"}", o.Name))
				req, err := http.NewRequest("GET", fmt.Sprintf("%s/vmspec", o.RootUrl), body)
				if err != nil {
					return err
				}
				req.Header.Add("Content-Type", "application/json")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				bytes, _ := ioutil.ReadAll(resp.Body)
				o.Output.Write(bytes)

				return nil
			}())
		},
	})

	return cmds
}

const (
	tplConnection = `{
"ConfigName"     : "{{ .Name }}",
"ProviderName"   : "{{ .CSP | ToUpper }}", 
"DriverName"     : "{{ .CSP }}-driver-v1.0", 
"CredentialName" : "{{ .Credential }}", 
"RegionName"     : "{{ .Region }}"
}`
)
