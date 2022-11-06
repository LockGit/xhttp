package process

import (
	"encoding/json"
	"sync"
	"xhttp/handler"
)

type IProcess interface {
	Name() string
	Exec(ctx *handler.Context)
}

var DefaultProcess = map[string]IProcess{}

func Register(p IProcess) {
	if _, ok := DefaultProcess[p.Name()]; !ok {
		DefaultProcess[p.Name()] = p
	}
}

type Response struct {
	Data map[string]interface{}
	mu   sync.Mutex
}

func NewResponse() *Response {
	return &Response{
		Data: make(map[string]interface{}),
	}
}

func (r *Response) Set(k string, v interface{}) {
	r.mu.Lock()
	r.Data[k] = v
	r.mu.Unlock()
}

func (r *Response) Get(k string) (v interface{}, ok bool) {
	v, ok = r.Data[k]
	return
}

func (r *Response) ToBytes() (bs []byte) {
	bs, _ = json.MarshalIndent(r.Data, "", "  ")
	return
}
