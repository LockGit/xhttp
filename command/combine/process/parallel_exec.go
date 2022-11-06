package process

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"xhttp/command/combine/process/params"
	"xhttp/handler"
	"xhttp/storage"
)

func init() {
	Register(NewParallelProcess())
}

type ParallelProcess struct {
}

func NewParallelProcess() *ParallelProcess {
	return &ParallelProcess{}
}

func (p *ParallelProcess) Name() string {
	return "parallel"
}

func (p *ParallelProcess) Exec(ctx *handler.Context) {
	data := NewResponse()
	wg := sync.WaitGroup{}
	for _, apiChild := range ctx.GetCurAPIChildren() {
		log.Println("start exec api:", apiChild.Name)
		requestParams := make(map[string]interface{})
		for _, param := range apiChild.Params {
			v, _ := params.GetParamsValue(ctx, param)
			requestParams[param.Name] = v
		}
		wg.Add(1)
		go func(api *storage.APIChildren) {
			if err := execRequest(api, requestParams, data); err != nil {
				log.Println(fmt.Sprintf("url:%s,child:%s,request err:%s", ctx.API.Url, apiChild.Name, err.Error()))
			}
			wg.Done()
		}(apiChild)
	}
	wg.Wait()

	ctx.Response.Header().Add("Content-Type", "application/json")
	ctx.Response.WriteHeader(http.StatusOK)
	ctx.Response.Write(data.ToBytes())
}
