package handler

import (
	"net/http"
	"xhttp/storage"
)

type Context struct {
	Request          *http.Request
	Response         http.ResponseWriter
	API              *storage.API
	MatchRouteParams map[string]string
}

func (c *Context) GetCurAPI() *storage.API {
	return c.API
}

func (c *Context) GetCurAPIChildren() (children []*storage.APIChildren) {
	if c.API == nil {
		return
	}
	return c.API.Children
}
