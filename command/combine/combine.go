package combine

import (
	"net/http"
	"xhttp/command"
	"xhttp/command/combine/process"
	"xhttp/handler"
)

func init() {
	command.Register("combine", setup)
}

func setup(nextHandler handler.Handler) (h handler.Handler, err error) {
	c := &Combine{
		NextHandler: nextHandler,
	}
	return c, nil
}

type Combine struct {
	NextHandler handler.Handler
}

func (c *Combine) ServerHTTP(ctx *handler.Context) {
	if c.Next() != nil {
		defer c.Next().ServerHTTP(ctx)
	}
	p, ok := process.DefaultProcess[ctx.API.ExecType]
	if !ok {
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		ctx.Response.Write([]byte("no support exec method"))
		return
	}
	p.Exec(ctx)
}

func (c *Combine) Next() handler.Handler {
	return c.NextHandler
}
