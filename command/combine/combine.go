package combine

import (
	"log"
	"net/http"
	"xhttp/command"
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
	if c.NextHandler != nil {
		defer c.NextHandler.ServerHTTP(ctx)
	}
	log.Println("Combine:", ctx.Request.RequestURI)
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write([]byte("ok:" + ctx.Request.RequestURI))
}

func (c *Combine) Next() handler.Handler {
	return c.NextHandler
}
