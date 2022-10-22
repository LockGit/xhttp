package handler

type Handler interface {
	ServerHTTP(ctx *Context)
	Next() Handler
}

type HandlerFunc func(ctx *Context)

func (f HandlerFunc) ServerHTTP(ctx *Context) {
	f(ctx)
}

func (f HandlerFunc) Next() Handler {
	return nil
}
