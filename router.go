package minigin

type Response struct {
	handlers HandlerChain
	params   Params
}

type Tree interface {
	Insert(method string, path string, handlers HandlerChain) error
	Search(method string, path string) (*Response, error)
}

type Router struct {
	tree Tree
}

var NilHandler = func(c *Context) {
	c.Writer.WriteHeader(404)
}

func NewRouter(tree Tree) *Router {
	return &Router{
		tree: tree,
	}
}

func (r Router) AddRoute(method, path string, handlers HandlerChain) {
	r.tree.Insert(method, path, handlers)
}

func (r Router) Search(method, path string) (*Response, error) {
	return r.tree.Search(method, path)
}

//TODO,trieを純粋にpath検索のロジックのみにして、routerにparam処理のロジックを移す
func (r Router) handle(c *Context) {
	res, err := r.Search(c.Request.Method, c.Request.URL.Path)
	if err != nil {
		//NilHandlerをセット
		return
	}
	//ここでcontextのparamsにセット
	c.Params = res.params
	c.handlers = append(c.handlers, res.handlers...)
	c.Next()
	return

}
