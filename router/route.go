package router

import (
	"log"
	"strings"
	"xhttp/command"
	"xhttp/handler"
	"xhttp/storage"
)

type Route struct {
	*storage.Project
	roots map[string]*node
}

func NewRoute(p *storage.Project) *Route {
	return &Route{
		Project: p,
		roots:   make(map[string]*node),
	}
}

func (r *Route) BuildTire(apis []*storage.API) {
	for _, api := range apis {
		log.Println("build tire url:", api.Url)
		method := strings.ToLower(api.Method)
		r.addRoute(method, api)
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *Route) addRoute(method string, api *storage.API) {
	pattern := api.Url
	parts := parsePattern(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0, api)
}

func (r *Route) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *Route) Match(ctx *handler.Context, url string) (api *storage.API, params map[string]string) {
	method := strings.ToLower(ctx.Request.Method)
	p := make(map[string]string)
	var n *node
	n, p = r.getRoute(method, url)
	if n == nil {
		return
	}
	return n.val, p
}

func (r *Route) ServerHTTP(ctx *handler.Context) {
	if r.Next() != nil {
		defer r.Next().ServerHTTP(ctx)
	}
	r.ExecCmds(ctx)
}

func (r *Route) ExecCmds(ctx *handler.Context) {
	command.GetCmdExecutor().ServerHTTP(ctx)
}

func (r *Route) Next() handler.Handler {
	return nil
}
