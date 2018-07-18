package impl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	hc "github.com/ddliu/go-httpclient"
)

const (
	Get    = iota //0
	Post          //1
	Put           //2
	Patch         //3
	Delete        //4
)

type HttpClient struct {
	hc *hc.HttpClient
}

func (h *HttpClient) Init(proxy string) (e error) {
	h.hc = hc.NewHttpClient()

	if proxy != "" {
		m := make(map[interface{}]interface{})
		m[hc.OPT_PROXY] = proxy
		h.hc.Defaults(m)
	}

	return nil
}

func (h *HttpClient) GetHttpClient() *hc.HttpClient {
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

func (h *HttpClient) DoRequest(httpUrl string, params map[string]string, httpMethod int, header map[string]string) (ret []byte, e error) {
	ret, _, e = h.DoRequest2(httpUrl, params, httpMethod, header)
	return
}

func (h *HttpClient) DoRequest2(httpUrl string, params map[string]string, httpMethod int, header map[string]string) (ret []byte, statsCode int, e error) {
	ret, statsCode, _, e = h.DoRequest2WithRetHeader(httpUrl, params, httpMethod, header)
	return
}
func (h *HttpClient) DoRequest2WithRetHeader(httpUrl string, params map[string]string, httpMethod int, header map[string]string) (ret []byte, statsCode int, rspHeader http.Header, e error) {

	if params == nil {
		params = make(map[string]string)
	}

	if header != nil && len(header) > 0 {
		h.hc.WithHeaders(header)
	}

	var res *hc.Response

	switch httpMethod {
	case Get:
		res, e = h.hc.Get(httpUrl, params)
	case Post:
		res, e = h.hc.Post(httpUrl, params)
	case Put:
		res, e = h.hc.PutJson(httpUrl, &params)
	case Patch:
		res, e = h.hc.Patch(httpUrl, params)
	case Delete:
		res, e = h.hc.Delete(httpUrl, params)
	}

	if e != nil {
		return
	}

	defer res.Body.Close()

	rspHeader = res.Header
	statsCode = res.StatusCode

	switch statsCode {
	case 200, 301, 302:
		ret, e = res.ReadAll()
	default:
		ret, _ = res.ReadAll()
		e = fmt.Errorf("http error statusCode=%v", res.StatusCode)
	}

	return
}

func (h *HttpClient) DoRequestPostJson(httpUrl string, jsonData interface{}, header map[string]string) (ret []byte, statsCode int, e error) {
	ret, statsCode, _, e = h.DoRequestPostJsonWithRetHeader(httpUrl, jsonData, header)
	return
}

func (h *HttpClient) DoRequestPostJsonWithRetHeader(httpUrl string, jsonData interface{}, header map[string]string) (ret []byte, statsCode int, rspHeader http.Header, e error) {
	if jsonData == nil {
		jsonData = make(map[string]string)
	}

	if header != nil && len(header) > 0 {
		h.hc.WithHeaders(header)
	}

	var res *hc.Response

	res, e = h.hc.PostJson(httpUrl, jsonData)

	if e != nil {
		return
	}

	defer res.Body.Close()

	rspHeader = res.Header

	statsCode = res.StatusCode

	switch statsCode {
	case 200, 301, 302:
		ret, e = res.ReadAll()
	default:
		ret, _ = res.ReadAll()
		e = fmt.Errorf("http error statusCode=%v", res.StatusCode)
	}

	return
}

func (h *HttpClient) DoRequestRetJson(httpUrl string, params map[string]string, httpMethod int, ret interface{}, header map[string]string) (e error) {
	var content []byte
	content, e = h.DoRequest(httpUrl, params, httpMethod, header)
	if e != nil {
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
