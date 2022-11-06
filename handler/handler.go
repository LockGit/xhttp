package handler

type Handler interface {
	ServerHTTP(ctx *Context)
	Next() Handler
}

type HandleFunc func(ctx *Context)

func (f HandleFunc) ServerHTTP(ctx *Context) {
	f(ctx)
}

func (f HandleFunc) Next() Handler {
	return nil
}
