package router

import (
	"log"
	"net"
	"xhttp/handler"
)

type Project struct {
	Projects map[string]*Route
	//host router
	//[project]Project
}

func (p *Project) Match(ctx *handler.Context) (node handler.Handler) {
	hostPort := ctx.Request.Header.Get("Host")
	host, _, _ := net.SplitHostPort(hostPort)
	node = p.Projects[host]
	log.Println("node is:", node)
	return
}
