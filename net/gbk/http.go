package gbk

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/axgle/mahonia"
)

var gbk mahonia.Decoder
var regCharset *regexp.Regexp

func HttpGetEx(url string, vparams ...map[string]string) (retBytes []byte, retE error) {

	if url == "" {
		return nil, errors.New("url is empty")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	retBytes, retE = ioutil.ReadAll(resp.Body)

	if retE != nil {
		return
	}

	var params map[string]string
	if len(vparams) > 0 {
		params = vparams[0]
	}

	for f := true; f; f = false {
		charset := ""

		contentType := strings.ToLower(resp.Header.Get("Content-Type"))

		if regCharset == nil {
			regCharset = regexp.MustCompile(`charset=([a-z0-9]+)`)
		}
		regRst := regCharset.FindStringSubmatch(contentType)
		if len(regRst) > 1 {
			charset = regRst[1]
		}

		if charset == "" && params != nil {
			charset, _ = params["charset"]
		}

		charset = strings.ToLower(charset)
		if charset == "gbk" || charset == "gb2312" {
			if gbk == nil {
				gbk = mahonia.NewDecoder("gbk")
			}
			retBytes = []byte(gbk.ConvertString(string(retBytes)))
		}
	}

	return retBytes, retE
}
