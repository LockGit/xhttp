package router

import (
	"log"
	"net"
	"xhttp/handler"
	"xhttp/storage"
)

type Projects struct {
	ProjectsMap map[string]*Route
}

func (p *Projects) Match(ctx *handler.Context) (route *Route, api *storage.API) {
	hostPort := ctx.Request.Header.Get("Host")
	host, _, _ := net.SplitHostPort(hostPort)
	log.Println("host is:", host)
	var ok bool
	host = "hello" //@todo change
	if route, ok = p.ProjectsMap[host]; !ok {
		return
	}
	log.Println("route is:", route.APIs[0].Url)
	//@todo url match 路由匹配实现
	return route, route.APIs[0]
}
