package router

import (
	"xhttp/command"
	"xhttp/handler"
)

type Route struct {
	ProjectName string
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
