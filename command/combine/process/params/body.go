package params

import (
	"encoding/json"
	"strings"
	"xhttp/handler"
)

func init() {
	Register("$.body", &Body{})
}

type Body struct {
}

func (b *Body) Get(ctx *handler.Context, name string) (v string, has bool) {
	ct := ctx.Request.Header.Get("Content-Type")
	if strings.HasPrefix(strings.ToLower(ct), "application/x-www-form-urlencoded") {
		return ctx.Request.Form.Get(name), ctx.Request.Form.Has(name)
	}
	if strings.HasPrefix(strings.ToLower(ct), "application/json") {
		decoder := json.NewDecoder(ctx.Request.Body)
		var params map[string]string
		if err := decoder.Decode(&params); err == nil {
			v, has = params[name]
			return
		}
	}
	return
}
