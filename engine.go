package minigin

import (
	"net/http"
)

const defaultMultipartMemory = 32 << 20 // 32 MB

type Content interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Engine struct {
	allowMultipartMemory int64
	RouterGroup          *RouterGroup
	router               *Router
}

type HandlerFunc func(c *Context)

type HandlerChain []HandlerFunc

func New() *Engine {

	e := &Engine{
		allowMultipartMemory: defaultMultipartMemory,
	}

	tree := NewTrie()
	e.router = NewRouter(tree)
	e.RouterGroup = &RouterGroup{
		basePath: "/",
		Handlers: nil,
		engine:   e,
	}
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

func (e *Engine) GET(path string, handlers ...HandlerFunc) {
	e.RouterGroup.AddRoute(http.MethodGet, path, handlers)
}

func (e *Engine) POST(path string, handlers ...HandlerFunc) {
	e.RouterGroup.AddRoute(http.MethodPost, path, handlers)
}

func (e *Engine) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return e.RouterGroup.Group(relativePath, handlers...)
}
