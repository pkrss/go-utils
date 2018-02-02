package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	pkBeans "github.com/pkrss/go-utils/beans"
)

type ControllerInterface interface {
	// Get()
	// Post()
	// Delete()
	// Put()
	// Patch()

	CloneAttribute(src ControllerInterface)

	OnPrepare()
	OnLeave()

	SetUrlParameters(p map[string]string)
	SetResponseWriter(w http.ResponseWriter)
	GetResponseWriter() http.ResponseWriter
	SetRequest(r *http.Request)
	GetRequest() *http.Request
	RequestBodyToJsonObject(ob interface{}) error

	AjaxError(message string, codes ...int)

	Header(key string) string
	SetCookieValue(key string, val string, maxAgeSecondss ...int)
	CookieValue(key string) string
	GetInt(key string, defValues ...int) int
	GetId64(key string, defValues ...int64) int64
	GetString(key string, defValues ...string) string
}

type Controller struct {
	W http.ResponseWriter
	R *http.Request

	UrlParameters map[string]string
}

func (this *Controller) OnPrepare() {

}

func (this *Controller) OnLeave() {

}

func (this *Controller) CloneAttribute(src ControllerInterface) {

}

// func (this *Controller) Get() {
// 	this.AjaxError("Not implement!")
// }
// func (this *Controller) Post() {
// 	this.AjaxError("Not implement!")
// }
// func (this *Controller) Delete() {
// 	this.AjaxError("Not implement!")
// }
// func (this *Controller) Put() {
// 	this.AjaxError("Not implement!")
// }
// func (this *Controller) Patch() {
// 	this.Patch()
// }

func (this *Controller) SetUrlParameters(p map[string]string) {
	this.UrlParameters = p
}

func (this *Controller) SetResponseWriter(w http.ResponseWriter) {
	this.W = w
}

func (this *Controller) SetRequest(r *http.Request) {
	this.R = r
}

func (this *Controller) GetResponseWriter() http.ResponseWriter {
	return this.W
}
func (this *Controller) GetRequest() *http.Request {
	decoder := json.NewDecoder(this.R.Body)
	err := decoder.Decode(&t)
	return this.R
}
func (this *Controller) RequestBodyToJsonObject(ob interface{}) error {
	decoder := json.NewDecoder(this.R.Body)
	err := decoder.Decode(ob)
	this.R.Body.Close()
	return err
}

/*
	\param key string, ":k": search in url path paramters, else search in url paramters
*/
func (this *Controller) Query(key string) string {
	if key == "" {
		return ""
	}

	if key[0] == ':' {
		if len(key) == 1 {
			return ""
		}
		v, _ := this.UrlParameters[key[1:]]
		return v
	}

	this.R.ParseForm()
	return this.R.Form.Get(key)
}

// Header returns request header item string by a given string.
// if non-existed, return empty string.
func (this *Controller) Header(key string) string {
	return this.R.Header.Get(key)
}

// Cookie returns request cookie item string by a given key.
// if non-existed, return empty string.
func (this *Controller) Cookie(key string) *http.Cookie {
	c, e := this.R.Cookie(key)
	if e != nil {
		return nil
	}
	return c
}

func (this *Controller) CookieValue(key string) string {
	c := this.Cookie(key)
	if c == nil {
		return ""
	}
	return c.Value
}

func (this *Controller) SetCookieValue(key string, val string, maxAgeSecondss ...int) {
	if val == "" {
		c := this.Cookie(key)
		if c != nil {
			c.Expires = time.Unix(0, 0)
		}
		return
	}
	maxAgeSecond := 0
	if len(maxAgeSecondss) > 0 {
		maxAgeSecond = maxAgeSecondss[0]
	}
	expiration := time.Now().Add(maxAgeSecond * time.Second)
	cookie := http.Cookie{Name: key, Value: val, Path: "/", Expires: expiration, MaxAge: maxAgeSecond}
	http.SetCookie(this.W, &cookie)
}

func (this *Controller) JsonResult(out interface{}) {
	this.W.Header().Add("Content-Type", "application/json; charset=utf-8")

	content, err := json.Marshal(out)

	if err != nil {
		http.Error(this.W, err.Error(), http.StatusInternalServerError)
		return
	}

	this.W.Write(content)
}

func (this *Controller) AjaxUnAuthorized(message string, codes ...int) {
	this.W.WriteHeader(401)
	this.AjaxError(message, codes...)
}

func (this *Controller) AjaxError(message string, codes ...int) {
	out := make(map[string]interface{})

	code := -1
	if len(codes) > 0 {
		code = codes[0]
	}

	out["code"] = code
	out["data"] = nil

	out["message"] = message

	this.JsonResult(out)
}

func (this *Controller) AjaxMsg(code int, data interface{}, message ...string) {
	out := make(map[string]interface{})
	out["code"] = code
	out["data"] = data

	var msg string

	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = "Operator succeed!"
	}

	out["message"] = msg

	this.JsonResult(out)
}

func (this *Controller) GetPageableFromRequest() *pkBeans.Pageable {
	var pageable pkBeans.Pageable
	var limit, pageNumber int

	limit = this.GetInt("pageSize", -1)
	if limit == -1 {
		limit = this.GetInt("limit", -1)
	}
	if limit == -1 {
		limit = 20
	}
	pageable.PageSize = limit

	pageNumber = this.GetInt("pageNumber", -1)
	if pageNumber == -1 {
		offset := this.GetInt("offset", -1)
		if offset == -1 {
			offset = 0
		} else {
			pageable.RspCodeFormat = true
			pageable.OffsetOldField = offset
		}
		if limit > 0 {
			pageNumber = 1 + (offset+limit-1)/limit
		}
	}
	if pageNumber < 1 {
		pageNumber = 1
	}

	pageable.PageNumber = pageNumber

	condArr := make(map[string]string, 0)

	query_params := []string{"q", "in_name", "in_list", "in_type", "like_name", "like_value"}
	for _, param := range query_params {
		paramValue := this.GetString(param)
		if len(paramValue) > 0 {
			condArr[param] = paramValue
		}
	}

	pageable.CondArr = condArr

	sort := this.GetString("sort")
	if pageable.Sort == "" {
		order := this.GetString("order")
		if order != "" {
			order = strings.ToLower(order)
			orders := strings.Split(order, " ")
			if len(orders) >= 2 {
				if orders[1] == "asc" {
					sort = "+" + orders[0]
				} else if orders[1] == "desc" {
					sort = "-" + orders[0]
				}
			}
		}
	}
	pageable.Sort = sort

	return &pageable
}

func (this *Controller) AjaxDbRecord(ob interface{}, oldCodeFormat ...bool) {
	if len(oldCodeFormat) > 0 && oldCodeFormat[0] {
		this.AjaxMsg(0, ob)
	} else {
		this.JsonResult(ob)
	}
}

func (this *Controller) AjaxDbList(pageable *pkBeans.Pageable, list interface{}, listSize int, total int64, oldCodeFormat ...bool) {

	var page pkBeans.Page
	if pageable.PageNumber == 1 {
		page.First = true
	} else {
		page.First = false
	}

	total2 := int(total)

	page.Number = pageable.PageNumber
	page.NumberOfElements = listSize
	page.Content = list
	page.TotalElements = total2
	page.Size = pageable.PageSize
	if page.Size > 0 {
		page.TotalPages = (total2 + page.Size - 1) / page.Size
	} else {
		page.TotalPages = 0
	}

	if pageable.PageNumber >= page.TotalPages-1 {
		page.Last = true
	} else {
		page.Last = false
	}

	rspCodeFormat := false
	if len(oldCodeFormat) > 0 {
		rspCodeFormat = oldCodeFormat[0]
	} else {
		rspCodeFormat = pageable.RspCodeFormat
	}

	if rspCodeFormat {
		rsp := make(map[string]interface{}, 0)
		rsp["list"] = page.Content
		rsp["listSize"] = page.NumberOfElements
		rsp["offset"] = (page.Number - 1) * page.Size
		rsp["limit"] = page.Size
		rsp["totalPages"] = page.TotalElements
		rsp["totalItemsCount"] = page.TotalElements
		rsp["links"] = nil

		this.AjaxMsg(0, rsp)
	} else {
		this.JsonResult(page)
	}
}

func (this *Controller) GetInt(key string, defValues ...int) int {
	defValue := 0
	if len(defValues) > 0 {
		defValue = defValues[0]
	}
	s := this.Query(key)
	if s == "" {
		return defValue
	}
	i, e := strconv.Atoi(s)
	if e != nil {
		return defValue
	}
	return i
}

func (this *Controller) GetId64(key string, defValues ...int64) int64 {
	var defValue int64 = 0
	if len(defValues) > 0 {
		defValue = defValues[0]
	}
	s := this.Query(key)
	if s == "" {
		return defValue
	}
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return defValue
	}
	return i
}

func (this *Controller) GetString(key string, defValues ...string) string {
	defValue := ""
	if len(defValues) > 0 {
		defValue = defValues[0]
	}
	s := this.Query(key)
	if s == "" {
		return defValue
	}

	return s
}
