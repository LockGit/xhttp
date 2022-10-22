package engine

import (
	"log"
	"net/http"
	"sync"
	"xhttp/handler"
	"xhttp/router"
)

type Engine struct {
	*router.Project
	once         sync.Once
	HookFuncList []func() error
}

func NewEngine() *Engine {
	return &Engine{
		Project: &router.Project{},
	}
}

func (e *Engine) InitProject() {
	projects := make(map[string]*router.Route)
	projects[""] = &router.Route{
		ProjectName: "hello",
	}
	e.Projects = projects
}

func (e *Engine) RegisterHook(hook ...func() error) {
	e.HookFuncList = append(e.HookFuncList, hook...)
}

func (e *Engine) RunHook() {
	for _, fn := range e.HookFuncList {
		go func(f func() error) {
			if err := f(); err != nil {
				log.Println("RunHook() err:", err.Error())
			}
		}(fn)
	}
}

func (e *Engine) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	e.once.Do(func() {
		e.RunHook()
		e.InitProject()
	})
	ctx := &handler.Context{
		Request:  request,
		Response: response,
	}
	node := e.Project.Match(ctx)
	if node != nil {
		node.ServerHTTP(ctx)
		return
	}
	response.WriteHeader(http.StatusNotFound)
	_, err := response.Write([]byte("404 not found"))
	if err != nil {
		log.Println("response.Write err:", err.Error())
	}
}
