package params

import (
	"xhttp/handler"
	"xhttp/storage"
)

var ParamsFactory = map[string]IParams{}

type IParams interface {
	Get(ctx *handler.Context, name string) (v string, has bool)
}

func Register(source string, p IParams) {
	if _, ok := ParamsFactory[source]; !ok {
		ParamsFactory[source] = p
	}
}

func GetParamsValue(ctx *handler.Context, param *storage.Param) (v string, hit bool) {
	if factory, ok := ParamsFactory[param.Source]; ok {
		val, has := factory.Get(ctx, param.Name)
		if !has {
			val = param.DefaultValue
		}
		return val, true
	}
	return param.DefaultValue, false
}
