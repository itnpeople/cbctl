package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type FilenamesFlags struct {
	filenames []string
}

func NewFilenamesFlags(filenames []string) *FilenamesFlags {
	return &FilenamesFlags{
		filenames: filenames,
	}
}
func (self *FilenamesFlags) ToBytes() ([]byte, error) {

	data := []byte{}
	for _, s := range self.filenames {

		var buf []byte
		var err error
		switch {
		case s == "-":
			buf, err = ioutil.ReadAll(os.Stdin)
		case strings.Index(s, "http://") == 0 || strings.Index(s, "https://") == 0:
			_, err = url.Parse(s)
			if err != nil {
				return nil, fmt.Errorf("the URL passed to filename %q is not valid: %v", s, err)
			}
			if resp, err := http.Get(s); err != nil {
				return nil, err
			} else {
				defer resp.Body.Close()
				buf, err = ioutil.ReadAll(resp.Body)
			}
		default:
			buf, err = ioutil.ReadFile(s)
		}
		if err != nil {
			return nil, err
		} else {
			data = append(data, buf...)
		}
	}

	return data, nil

}
