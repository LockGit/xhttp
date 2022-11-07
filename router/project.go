package router

import (
	"log"
	"xhttp/handler"
	"xhttp/storage"
)

type Projects struct {
	ProjectsMap map[string]*Route
}

func (p *Projects) Match(ctx *handler.Context) (route *Route, api *storage.API, params map[string]string) {
	project := ctx.Request.Header.Get("X-Project")
	log.Println("X-Project is:", project)
	var ok bool
	if route, ok = p.ProjectsMap[project]; !ok {
		return
	}
	mp := make(map[string]string)
	api, mp = route.Match(ctx, ctx.Request.RequestURI)
	return route, api, mp
}
