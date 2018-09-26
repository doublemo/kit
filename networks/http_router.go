package networks

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/doublemo/kit/networks/tree"
)

type Route struct {
	Method   string
	Path     string
	Params   url.Values
	handlers WebHandlersChain
}

// Router 路由存储
type Router struct {
	tree *tree.Node
}

// AddRoute 增加路由
func (r *Router) AddRoute(method, path string, handlers WebHandlersChain) error {
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("Prefix path invalid :%s %s", method, path)
	}
	route := &Route{
		Method:   method,
		Path:     path,
		Params:   make(url.Values),
		handlers: handlers,
	}
	return r.tree.Add(r.TreePath(method, path), route)
}

func (r *Router) Find(method, path string) *Route {
	leaf, expansions := r.tree.Find(r.TreePath(method, path))
	if leaf == nil {
		return nil
	}

	route, ok := leaf.Value.(*Route)
	if !ok {
		return nil
	}

	if len(expansions) > 0 {
		route.Params = make(url.Values)
		for i, v := range expansions {
			route.Params.Set(leaf.Wildcards[i], v)
		}
	}

	return route
}

func (r *Router) TreePath(method, path string) string {
	if method == "*" {
		method = ":METHOD"
	}

	return "/" + strings.ToUpper(method) + path
}

func NewRouter() *Router {
	return &Router{tree: tree.New()}
}
