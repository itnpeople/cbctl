package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
)

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}

func ToPrettyJSON(data []byte) []byte {

	if len(data) > 0 {
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", "  "); err == nil {
			return buf.Bytes()
		}
	}
	return data
}

func ToTemplateBytes(tpl string, todo interface{}) ([]byte, error) {

	t, err := template.New("tpl").Funcs(
		template.FuncMap{
			"ToUpper": strings.ToUpper,
		}).Parse(tpl)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = t.Execute(&out, todo)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil

}

func ToFilenameBytes(filename string) ([]byte, error) {

	var buf []byte
	var err error
	switch {
	case filename == "-":
		buf, err = ioutil.ReadAll(os.Stdin)
	case strings.Index(filename, "http://") == 0 || strings.Index(filename, "https://") == 0:
		if _, err = url.Parse(filename); err == nil {
			if resp, err := http.Get(filename); err == nil {
				defer resp.Body.Close()
				buf, err = ioutil.ReadAll(resp.Body)
			}
		}
	default:
		buf, err = ioutil.ReadFile(filename)
	}

	if err != nil {
		return nil, err
	}

	//var body interface{}
	if buf, err = yaml.YAMLToJSON(buf); err != nil {
		return nil, err
	}

	return buf, err

}

func GetKeys(data map[string]interface{}) []interface{} {

	var keys []interface{}

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			keys = append(keys, k.Interface())
		}

	}

	return keys

}
