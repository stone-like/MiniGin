package minigin

import (
	"net/http"
	"path/filepath"
)

type RouterGroup struct {
	basePath string
	Handlers HandlerChain
	engine   *Engine
}

func NewRouterGroup(path string, handlers HandlerChain) *RouterGroup {
	return &RouterGroup{
		basePath: path,
		Handlers: handlers,
	}
}

func (r *RouterGroup) combineHandlers(handlers HandlerChain) HandlerChain {
	targetLen := len(r.Handlers) + len(handlers)
	combinedHandlers := make(HandlerChain, targetLen)
	copy(combinedHandlers, r.Handlers)
	copy(combinedHandlers[len(r.Handlers):], handlers)
	return combinedHandlers
}

//実際にユーザーが引数として渡す関数だけ ...HandlerFuncと可変長引数にする、そうすれば0個～複数の引数が取れて便利
func (r *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodGet, relativePath, handlers)
}

func (r *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	r.AddRoute(http.MethodPost, relativePath, handlers)
}

func (r *RouterGroup) AddRoute(method, relativePath string, handlers HandlerChain) {

	absolutePath := createAbsolutePath(r.basePath, relativePath)
	combinedHandlers := r.combineHandlers(handlers)
	r.engine.router.AddRoute(method, absolutePath, combinedHandlers)
}

func (r *RouterGroup) Use(handlers ...HandlerFunc) {
	r.Handlers = append(r.Handlers, handlers...)
}

//windowsではfileSeparatorが"\"なので、filepath.ToSlashで"/"に変換、別に物理パスじゃないのでpath.Joinでもいいけど
func createAbsolutePath(basePath, relativePath string) string {

	return filepath.ToSlash(filepath.Join(basePath, relativePath))
}

func (r *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	absolutePath := createAbsolutePath(r.basePath, relativePath)
	combinedHandlers := r.combineHandlers(handlers)
	return &RouterGroup{
		basePath: absolutePath,
		Handlers: combinedHandlers,
		engine:   r.engine,
	}
}
