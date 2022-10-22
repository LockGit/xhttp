package command

import (
	"xhttp/handler"
)

type Setup func(cmd handler.Handler) (h handler.Handler, err error)

var Cmds map[string]Setup

func init() {
	Cmds = make(map[string]Setup)
}

func Register(name string, setup Setup) {
	if _, ok := Cmds[name]; !ok {
		Cmds[name] = setup
	}
}

func Get(name string) (setup Setup) {
	if v, ok := Cmds[name]; ok {
		return v
	}
	return
}
