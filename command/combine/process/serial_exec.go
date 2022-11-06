package process

import (
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"strings"
	"xhttp/command/combine/process/params"
	"xhttp/handler"
)

func init() {
	Register(NewSerialProcess())
}

type SerialProcess struct {
}

func NewSerialProcess() *SerialProcess {
	return &SerialProcess{}
}

func (s *SerialProcess) Name() string {
	return "serial"
}

func (s *SerialProcess) Exec(ctx *handler.Context) {
	data := NewResponse()
	for _, api := range ctx.GetCurAPIChildren() {
		requestParams := make(map[string]interface{})
		for _, param := range api.Params {
			v, hit := params.GetParamsValue(ctx, param)
			if hit {
				requestParams[param.Name] = v
			} else {
				//get value from other api result
				arr := strings.SplitN(param.Source, ".", 3)
				if len(arr) == 2 {
					if val, ok := data.Get(arr[1]); ok {
						requestParams[param.Name] = val
					}
				} else if len(arr) == 3 {
					if val, ok := data.Get(arr[1]); ok {
						key := arr[2]
						value := gjson.Get(val.(string), key)
						requestParams[param.Name] = value.String()
					}
				} else {
					requestParams[param.Name] = ""
				}
			}
		}
		if err := execRequest(api, requestParams, data); err != nil {
			log.Println(fmt.Sprintf("url:%s,child:%s,request err:%s", ctx.API.Url, api.Name, err.Error()))
		}
	}
	ctx.Response.Header().Add("Content-Type", "application/json")
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(data.ToBytes())
}
