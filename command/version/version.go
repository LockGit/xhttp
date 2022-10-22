package version

import (
	"xhttp/command"
	"xhttp/handler"
)

func init() {
	command.Register("version", setup)
}

func setup(nextHandler handler.Handler) (h handler.Handler, err error) {
	v := &Version{
		NextHandler: nextHandler,
	}
	return v, nil
}

type Version struct {
	NextHandler handler.Handler
}

func (v *Version) ServerHTTP(ctx *handler.Context) {
	if v.NextHandler != nil {
		defer v.NextHandler.ServerHTTP(ctx)
	}
	ctx.Response.Header().Add("X-Version", "1.0.0")
}

func (v *Version) Next() handler.Handler {
	return v.NextHandler
}
