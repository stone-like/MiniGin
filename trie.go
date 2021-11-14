package minigin

import (
	"fmt"
	"strings"
)

func pathNotRegisteredError(path string) error {
	return fmt.Errorf("path %v is not registered", path)
}

func methodNotRegisteredError(method string) error {
	return fmt.Errorf("method %v is not registered", method)
}

const (
	pathRoot      string = "/"
	pathDelimiter string = "/"
)

type tree struct {
	node *node
}

type node struct {
	path      string
	children  map[string]*node
	actions   map[string]HandlerChain
	endOfPath bool
	isWild    bool
}

func NewTrie() *tree {
	return &tree{
		node: &node{
			path:     pathRoot,
			actions:  make(map[string]HandlerChain),
			children: make(map[string]*node),
		},
	}
}

func isWild(path string) bool {
	return path[0] == ':'
}

func getWildKey(path string) string {
	return path[1:]
}

func (t *tree) Insert(method string, path string, handlers HandlerChain) error {
	curNode := t.node

	//"/"の場合は直ぐに終了
	if path == pathRoot {
		curNode.path = path
		curNode.actions[method] = handlers
		curNode.endOfPath = true
		return nil
	}

	paths := explodePath(path)

	for i, path := range paths {
		if nextNode, ok := curNode.children[path]; ok {
			curNode = nextNode
			continue
		}

		curNode.children[path] = &node{
			path:     path,
			actions:  make(map[string]HandlerChain),
			children: make(map[string]*node),
			isWild:   isWild(path),
		}

		curNode = curNode.children[path]

		if i == len(paths)-1 {
			curNode.path = path
			curNode.actions[method] = handlers
			curNode.endOfPath = true
		}
	}

	return nil

}

func findChild(node *node, path string, fullPath string) (*node, bool, error) {
	if nextNode, ok := node.children[path]; ok {
		return nextNode, false, nil
	}

	for pathName, node := range node.children {
		if isWild(pathName) {
			return node, true, nil
		}
	}

	return nil, false, pathNotRegisteredError(fullPath)
}

//TODO trieは純粋にpathMatchingだけにして、methodやhandlerの登録情報はrouterに持たせた方がよさそう?

//完全一致しなければ返さないようにしている
func (t *tree) getPath(paths []string, targetPath string) (*node, Params, error) {

	var params []Param
	curNode := t.node

	if targetPath == "/" {
		return curNode, []Param{}, nil
	}

	for _, p := range paths {

		nextNode, isWild, err := findChild(curNode, p, targetPath)

		if err != nil {
			return nil, []Param{}, err
		}

		//:mameみたいなやつ処理をしなければいけない
		if isWild {
			params = append(params, Param{
				Key:   getWildKey(nextNode.path),
				Value: p,
			})
		}

		curNode = nextNode

	}

	//これは例えば登録されているのが/foo/bar/bazで/foo/barを検索したときに、
	//うまくここまで来るが、ただ検索Pathと完全一致でないためNGを返したい
	//なので登録したPathの最後,bazにはendOfPathがついていて、barにはendOfPathがついていない、これを使う
	if !curNode.endOfPath {
		return nil, []Param{}, pathNotRegisteredError(targetPath)
	}

	return curNode, params, nil

}

func (t *tree) Search(method string, path string) (*Response, error) {
	paths := explodePath(path)

	pathNode, params, err := t.getPath(paths, path)
	if err != nil {
		return nil, err
	}

	handlers, ok := pathNode.actions[method]
	if !ok {
		return nil, methodNotRegisteredError(method)
	}

	return &Response{
		handlers: handlers,
		params:   params,
	}, nil
}

func explodePath(path string) []string {
	paths := strings.Split(path, pathDelimiter)
	var list []string

	for _, p := range paths {
		if p == "" {
			continue
		}
		list = append(list, p)
	}

	return list
}
