package minigin

import (
	"net/http"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

type Content interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Engine struct {
	router               Router
	allowMultipartMemory int64
}

type HandlerFunc func(c *Context)

func New() *Engine {

	e := &Engine{
		allowMultipartMemory: defaultMultipartMemory,
	}
	// e.router = make(Router)
	return e
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := e.allocateContext()
	c.Request = r
	c.Writer = w

	//ここでmiddlewareを作り、handlerに設定

	e.handleHttpRequest(c)

}

func (e *Engine) handleHttpRequest(c *Context) {
	e.router.handle(c)
}

func (e *Engine) allocateContext() *Context {
	return &Context{engine: e}
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.router.AddRoute("GET", path, handler)
}

func (e *Engine) POST(path string, handler HandlerFunc) {
	e.router.AddRoute("POST", path, handler)
}
