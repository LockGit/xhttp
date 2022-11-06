package params

import "xhttp/handler"

func init() {
	Register("$.cookie", &Cookie{})
}

type Cookie struct {
}

func (c *Cookie) Get(ctx *handler.Context, name string) (v string, has bool) {
	cookies := ctx.Request.Cookies()
	for _, cookie := range cookies {
		if name == cookie.Name {
			return cookie.Value, true
		}
	}
	return v, false
}
