package http

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Method ...
type Method int

const (
	// MethodUnknown ...
	MethodUnknown Method = iota
	Get
	Post
	Put
	Patch
	Delete
	PutEmpty
	GetJSON
	PostJSON
	PutJSON
	PatchJSON
	DeleteJSON
)

func method2String(method Method) string {
	ret := ""
	switch method {
	case Get, GetJSON:
		ret = "GET"
	case Post, PostJSON:
		ret = "POST"
	case Put, PutEmpty, PutJSON:
		ret = "PUT"
	case Patch, PatchJSON:
		ret = "PATCH"
	case Delete, DeleteJSON:
		ret = "DELETE"
	}
	return ret
}
func addUrlParams(url_ string, params url.Values) string {
	if len(params) == 0 {
		return url_
	}

	v := params.Encode()
	if v != "" {
		if strings.Contains(url_, "?") {
			url_ += "&"
		} else {
			url_ += "?"
		}
		url_ += v
	}

	return url_
}
func toUrlValues(v interface{}) (ret url.Values, e error) {
	switch t := v.(type) {
	case url.Values:
		ret = t
	case map[string][]string:
		ret = url.Values(t)
	case map[string]string:
		ret = make(url.Values)
		for k, v := range t {
			ret.Add(k, v)
		}
	case nil:
		ret = make(url.Values)
	default:
		e = errors.New("Invalid toUrlValues value")
	}
	return
}

// Read response body into a byte slice.
func readAll(res *http.Response) ([]byte, error) {
	var reader io.ReadCloser
	var err error
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return ioutil.ReadAll(reader)
	default:
		return ioutil.ReadAll(res.Body)
	}

}

type HttpClient struct {
	header    map[string]string
	hc        *http.Client
	transport *http.Transport
}

func (h *HttpClient) Init(proxy string) (e error) {

	h.header = make(map[string]string)
	m := h.header

	transport := &http.Transport{
		// Proxy: ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	h.hc = &http.Client{Transport: transport}

	// transport.Dial = func(network, addr string) (net.Conn, error) {
	// 	var conn net.Conn
	// 	var err error
	// 	if connectTimeoutMS > 0 {
	// 		conn, err = net.DialTimeout(network, addr, time.Duration(connectTimeoutMS)*time.Millisecond)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	} else {
	// 		conn, err = net.Dial(network, addr)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	if timeoutMS > 0 {
	// 		conn.SetDeadline(time.Now().Add(time.Duration(timeoutMS) * time.Millisecond))
	// 	}

	// 	return conn, nil
	// }

	if proxy != "" {
		if !strings.Contains(proxy, "://") {
			proxy = "http://" + proxy
		}

		if proxyUrl, err := url.Parse(proxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		} else {
			e = err
			return
		}

		m["Pragma"] = "no-cache"
		m["Cache-Control"] = "no-cache, must-revalidate"

		m["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"
		m["Accept-Encoding"] = "gzip"
	}

	return nil
}

func (c *HttpClient) SendRequest(httpUrl string, params interface{}, method Method, header map[string]string) (res *http.Response, e error) {

	h := make(map[string]string)
	for k, v := range c.header {
		if _, ok := h[k]; !ok {
			h[k] = v
		}
	}

	h["Content-Type"] = "application/x-www-form-urlencoded"

	var body io.Reader
	switch method {
	case Post:
		if params != nil {
			if paramsValues, e2 := toUrlValues(params); e2 == nil {
				body = strings.NewReader(paramsValues.Encode())
			} else {
				e = e2
				return
			}
		}
	case GetJSON, PostJSON, PutJSON, PatchJSON, DeleteJSON:
		h["Content-Type"] = "application/json"
		if params != nil {
			var content []byte
			switch t := params.(type) {
			case []byte:
				content = t
			case string:
				content = []byte(t)
			default:
				if content, e = json.Marshal(params); e != nil {
					return
				}
			}
			body = bytes.NewBuffer(content)
		}
	default:
		if paramsValues, e2 := toUrlValues(params); e2 == nil {
			httpUrl = addUrlParams(httpUrl, paramsValues)
		} else {
			e = e2
			return
		}
	}

	if header != nil {
		for k, v := range header {
			h[k] = v
		}
	}

	req, err := http.NewRequest(method2String(method), httpUrl, body)
	if err != nil {
		e = err
		return
	}

	for k, v := range h {
		req.Header.Set(k, v)
	}

	res, e2 := c.hc.Do(req)
	if e2 != nil {
		e = e2
		return
	}

	return
}

func (c *HttpClient) FetchResult(httpUrl string, params interface{}, method Method, header map[string]string) (ret []byte, statusCode int, rspHeader http.Header, e error) {

	res, e2 := c.SendRequest(httpUrl, params, method, header)

	if res != nil {
		rspHeader = res.Header
		statusCode = res.StatusCode
	}

	if e2 != nil {
		e = e2
		return
	}

	defer res.Body.Close()

	if ret, e = readAll(res); e != nil {
		return
	}

	if statusCode >= 400 {
		e = fmt.Errorf("http error statusCode=%d\n%s", statusCode, string(ret))

		// res.Body.Close()
	}

	return
}

// func (h *HttpClient) onDoCloudflareError(statsCode int, rspHeader http.Header) {
// 	if (statsCode >= 400) && (statsCode < 600) { // too busy  //
// 		if rspHeader == nil {
// 			return
// 		}
// 		v, ok := rspHeader["Server"]
// 		if ok && (len(v) > 0) && (strings.Contains(v[0], "cloudflare")) {
// 			time.Sleep(1 * time.Second)
// 			if OsExitWhenCloudflare50x {
// 				os.Exit(-97)
// 			}
// 			pkReflect.SetStructFieldValue(h.hc, "transport", nil, false)
// 		}
// 	}
// }

func (h *HttpClient) GetHttpClient() *http.Client {
	return h.hc
}

func (h *HttpClient) DoGet(httpUrl string, params map[string]string, header map[string]string) (ret []byte, e error) {
	return h.DoRequest(httpUrl, params, Get, header)
}
func (h *HttpClient) DoGetRetJson(httpUrl string, params map[string]string, rsp interface{}, header map[string]string) (e error) {
	return h.DoRequestRetJson(httpUrl, params, Get, rsp, header)
}

func (h *HttpClient) DoPost(httpUrl string, params map[string]string, header map[string]string) (ret []byte, e error) {
	return h.DoRequest(httpUrl, params, Post, header)
}
func (h *HttpClient) DoPostRetJson(httpUrl string, params map[string]string, rsp interface{}, header map[string]string) (e error) {
	return h.DoRequestRetJson(httpUrl, params, Post, rsp, header)
}

func (h *HttpClient) DoPut(httpUrl string, params map[string]string, header map[string]string) (ret []byte, e error) {
	return h.DoRequest(httpUrl, params, Put, header)
}
func (h *HttpClient) DoPutRetJson(httpUrl string, params map[string]string, rsp interface{}, header map[string]string) (e error) {
	return h.DoRequestRetJson(httpUrl, params, Put, rsp, header)
}
func (h *HttpClient) DoPutEmptyRetJson(httpUrl string, rsp interface{}, header map[string]string) (e error) {
	return h.DoRequestRetJson(httpUrl, nil, PutEmpty, rsp, header)
}

func (h *HttpClient) DoPatch(httpUrl string, params map[string]string, header map[string]string) (ret []byte, e error) {
	return h.DoRequest(httpUrl, params, Patch, header)
}
func (h *HttpClient) DoPatchRetJson(httpUrl string, params map[string]string, rsp interface{}, header map[string]string) (e error) {
	return h.DoRequestRetJson(httpUrl, params, Patch, rsp, header)
}

func (h *HttpClient) DoDelete(httpUrl string, params map[string]string, header map[string]string) (ret []byte, e error) {
	return h.DoRequest(httpUrl, params, Delete, header)
}
func (h *HttpClient) DoDeleteRetJson(httpUrl string, params map[string]string, rsp interface{}, header map[string]string) (e error) {
	return h.DoRequestRetJson(httpUrl, params, Delete, rsp, header)
}

func (h *HttpClient) DoRequest(httpUrl string, params interface{}, httpMethod Method, header map[string]string) (ret []byte, e error) {
	ret, _, e = h.DoRequest2(httpUrl, params, httpMethod, header)
	return
}

func (h *HttpClient) DoRequest2(httpUrl string, params interface{}, httpMethod Method, header map[string]string) (ret []byte, statsCode int, e error) {
	ret, statsCode, _, e = h.DoRequest2WithRetHeader(httpUrl, params, httpMethod, header)
	return
}
func (h *HttpClient) DoRequest2WithRetHeader(httpUrl string, params interface{}, httpMethod Method, header map[string]string) (ret []byte, statsCode int, rspHeader http.Header, e error) {

	return h.FetchResult(httpUrl, params, httpMethod, header)
}

func (h *HttpClient) DoRequest2WithRetRsp(httpUrl string, params interface{}, httpMethod Method, header map[string]string) (res *http.Response, e error) {
	return h.SendRequest(httpUrl, params, httpMethod, header)
}

// func (h *HttpClient) onDoCloudflareError(statsCode int, rspHeader http.Header) {
// 	if (statsCode >= 400) && (statsCode < 600) { // too busy  //
// 		if rspHeader == nil {
// 			return
// 		}
// 		v, ok := rspHeader["Server"]
// 		if ok && (len(v) > 0) && (strings.Contains(v[0], "cloudflare")) {
// 			time.Sleep(1 * time.Second)
// 			if OsExitWhenCloudflare50x {
// 				os.Exit(-97)
// 			}
// 			pkReflect.SetStructFieldValue(h.hc, "transport", nil, false)
// 		}
// 	}
// }

func (h *HttpClient) DoRequestPostJson(httpUrl string, jsonData interface{}, header map[string]string) (ret []byte, statsCode int, e error) {
	ret, statsCode, _, e = h.DoRequestPostJsonWithRetHeader(httpUrl, jsonData, header)
	return
}

func (h *HttpClient) DoRequestPostJsonWithRetHeader(httpUrl string, jsonData interface{}, header map[string]string) (ret []byte, statsCode int, rspHeader http.Header, e error) {
	return h.DoRequestJsonWithRetHeader2(httpUrl, jsonData, header, Post)
}

func (h *HttpClient) DoRequestJsonWithRetHeader2(httpUrl string, jsonData interface{}, header map[string]string, method Method) (ret []byte, statsCode int, rspHeader http.Header, e error) {
	switch method {
	case Get:
		method = GetJSON
	case Post:
		method = PostJSON
	case Put:
		method = PutJSON
	case Patch:
		method = PatchJSON
	case Delete:
		method = DeleteJSON
	}
	return h.FetchResult(httpUrl, jsonData, method, header)
}

func (h *HttpClient) SendJson(method Method, url string, data interface{}, headers map[string]string) (*http.Response, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"
	return h.SendRequest(url, data, method, headers)
}

func (h *HttpClient) DoRequestRetJson(httpUrl string, params map[string]string, httpMethod Method, ret interface{}, header map[string]string) (e error) {
	var content []byte

	content, _, _, e = h.FetchResult(httpUrl, params, httpMethod, header)
	if e != nil {
		log.Println(string(content))
		return
	}
	if ret != nil {
		return json.Unmarshal(content, ret)
	}
	return nil
}

func NewHttpClient(proxy string) (ret *HttpClient, e error) {
	ret = &HttpClient{}
	e = ret.Init(proxy)
	return
}

func HeaderAddBasicAuth(header map[string]string, username string, password string) map[string]string {
	if header == nil {
		header = make(map[string]string)
	}

	auth := username + ":" + password
	auth2 := base64.StdEncoding.EncodeToString([]byte(auth))

	header["Authorization"] = "Basic " + auth2
	return header
}
