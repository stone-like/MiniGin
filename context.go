package minigin

import (
	"fmt"
	"miniGin/binding"
	"net/http"
	"net/url"
)

const (
	MIMEJSON     = binding.MIMEJSON
	MIMEPOSTForm = binding.MIMEPOSTForm
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

type Context struct {
	engine     *Engine
	Request    *http.Request
	Writer     http.ResponseWriter
	Params     Params
	handlers   HandlerChain
	index      int8
	queryCache url.Values
	formCache  url.Values

	// CandidateValues url.Values
}

// func (c *Context) ParamToUrlValues() map[string][]string {
// 	m := make(map[string][]string)

// 	for _, param := range c.Params {
// 		m[param.Key] = []string{param.Value}
// 	}
// 	return m
// }

// func (c *Context) AddCandidateValue(values map[string][]string) {
// 	for key, value := range values {
// 		c.CandidateValues[key] = value
// 	}
// }

//名前付き返り値を使った方がこの場合GetQueryを使いまわせて少しスッキリして返せる(本当に少しだけど)
//基本関数は使いまわしてWrapperみたいにして使いたい
func (c *Context) Query(key string) (value string) {
	value, _ = c.GetQuery(key)
	return
}

func (c *Context) GetQuery(key string) (string, bool) {
	if values, ok := c.GetQueryArray(key); ok {
		return values[0], true
	}

	return "", false
}

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.Request != nil {
			c.queryCache = c.Request.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

func (c *Context) PostForm(key string) (value string) {
	value, _ = c.GetPostForm(key)
	return
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], true
	}

	return "", false
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		req := c.Request
		// req.ParseMultipartFormはinMemoryにparseするのでメモリを使う、どれくらい使うかを引数で指定
		if err := req.ParseMultipartForm(c.engine.allowMultipartMemory); err != nil {
			fmt.Println("error on parse multipart form array")
		}
		c.formCache = req.PostForm
	}
}

func (c *Context) GetPostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

//content-typeには; のあとにセットするような場合もあるので、filterする
//application/x-www-form-urlencoded; charset=utf-8
func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) ContentType() string {
	return filterFlags(c.requestHeader("Content-Type"))
}

func (c *Context) Bind(obj interface{}) error {
	b := binding.New(c.Request.Method, c.ContentType())
	return b.Bind(c.Request, obj)
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
