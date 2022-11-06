package params

import "xhttp/handler"

func init() {
	Register("$.query", &Query{})
}

type Query struct {
}

func (q Query) Get(ctx *handler.Context, name string) (v string, has bool) {
	v = ctx.Request.URL.Query().Get(name)
	return v, ctx.Request.URL.Query().Has(name)
}
