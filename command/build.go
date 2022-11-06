package command

import (
	"fmt"
	"log"
	"sync"
	"xhttp/handler"
)

var (
	cmdExecutor handler.Handler
	once        sync.Once
)
var CmdList = []string{
	"version",
	"combine",
}

func buildCmd() (h handler.Handler, e error) {
	exec := handler.Handler(handler.HandleFunc(func(ctx *handler.Context) {}))
	for i := len(CmdList) - 1; i >= 0; i-- {
		name := CmdList[i]
		if setup := Get(name); setup != nil {
			next, err := setup(exec)
			if err != nil {
				panic(fmt.Sprintf("buildCmd:%s err:%s", name, err.Error()))
			}
			exec = next
			log.Println("build cmd:", name, "===> ok")
		}
	}
	return exec, nil
}

func Build() (err error) {
	once.Do(func() {
		cmdExecutor, err = buildCmd()
		if err != nil {
			return
		}
	})
	return
}

func GetCmdExecutor() handler.Handler {
	return cmdExecutor
}
