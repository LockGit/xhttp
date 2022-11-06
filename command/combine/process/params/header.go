package params

import "xhttp/handler"

func init() {
	Register("$.header", &Header{})
}

type Header struct {
}

func (h *Header) Get(ctx *handler.Context, name string) (v string, has bool) {
	if len(ctx.Request.Header.Values(name)) > 0 {
		has = true
	}
	return ctx.Request.Header.Get(name), has
}
