package engine

import (
	"log"
	"net/http"
	"sync"
	"xhttp/handler"
	"xhttp/router"
	"xhttp/storage"
)

type Engine struct {
	store        *storage.Storage
	Projects     *router.Projects
	once         sync.Once
	HookFuncList []func() error
	lock         sync.Mutex
}

type Option func(engine *Engine)

func WithStorage(store *storage.Storage) Option {
	return func(engine *Engine) {
		engine.store = store
	}
}

func NewEngine(options ...Option) *Engine {
	engine := &Engine{
		Projects: nil,
	}
	for _, op := range options {
		op(engine)
	}
	return engine
}

func (e *Engine) updateProject(p *storage.Project) {
	e.lock.Lock()
	defer e.lock.Unlock()
	r := &router.Route{
		Project: p,
	}
	e.Projects.ProjectsMap[p.Name] = r
}

func (e *Engine) removeProject(p *storage.Project) {
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.Projects.ProjectsMap, p.Name)
}

func (e *Engine) InitProject() {
	projects, err := e.store.GetAll()
	if err != nil {
		log.Fatalf("Engine.InitProject() fatal,err:%s", err.Error())
	}
	ps := &router.Projects{
		ProjectsMap: make(map[string]*router.Route),
	}
	for _, p := range projects {
		if _, ok := ps.ProjectsMap[p.Name]; !ok {
			log.Println("add project: " + p.Name)
			r := &router.Route{
				Project: p,
			}
			ps.ProjectsMap[p.Name] = r
		}
	}
	e.Projects = ps

	go func() {
		if err = e.store.Watch(); err != nil {
			log.Fatalf("e.store.Watch() err:" + err.Error())
		}
	}()

	go func() {
		for ch := range e.store.WatchEvent() {
			log.Println("ch update", ch)
			switch ch.Op {
			case storage.OpMod, storage.OpAdd:
				log.Println("update project:", ch)
				e.updateProject(ch.Project)
			case storage.OpDel:
				log.Println("del project:", ch)
				e.removeProject(ch.Project)
			}
		}
	}()
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
	})
	ctx := &handler.Context{
		Request:  request,
		Response: response,
	}
	route, api := e.Projects.Match(ctx)
	ctx.API = api
	if route == nil {
		response.WriteHeader(http.StatusNotFound)
		_, err := response.Write([]byte("404 not found"))
		if err != nil {
			log.Println("response.Write err:", err.Error())
		}
		return
	}
	route.ServerHTTP(ctx)
}
