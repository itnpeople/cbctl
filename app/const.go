package app

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"

	"github.com/itnpeople/cbctl/utils"
)

/*
#!/usr/bin/env bash
version=0.1.0
go build -ldflags="-X 'github.com/itnpeople/cbctl/app.BuildTime=$(date -u +%FT%T%Z)' -X 'github.com/itnpeople/cbctl/app.BuildVersion=$version'" .
*/
var (
	BuildVersion string = ""
	BuildTime    string = ""
)

type IOStreams struct {
	In     *os.File
	Out    *os.File
	ErrOut *os.File
}

const (
	OUTPUT_JSON = "json"
	OUTPUT_YAML = "yaml"
)

type OutputType string

type Output struct {
	Type   *string
	Stream *os.File
}

func (self *Output) Write(json []byte) {
	if *self.Type == OUTPUT_JSON {
		self.Stream.Write(utils.ToPrettyJSON(json))
	} else {
		if d, err := yaml.JSONToYAML(json); err == nil {
			self.Stream.Write(d)
		} else {
			self.Stream.Write(json)
		}
	}
}

func (self *Output) WriteString(format string, params ...interface{}) {
	line := fmt.Sprintf(format, params...)
	self.Stream.WriteString(line)
}
