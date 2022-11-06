package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"xhttp/command"
	_ "xhttp/command/combine"
	_ "xhttp/command/version"
	"xhttp/engine"
	"xhttp/storage"
)

func main() {
	// init cmd
	if err := command.Build(); err != nil {
		log.Fatalf("command.Build() fatal error:%s", err.Error())
	}
	// use file storage
	fileStorage := storage.NewFileStorage()
	oStore := storage.NewStorage(fileStorage)
	if err := oStore.Init(); err != nil {
		log.Fatalf("oStore.Init() fatal error:%s", err.Error())
	}
	// init engine
	oEngine := engine.NewEngine(engine.WithStorage(oStore))
	oEngine.InitProject()
	oEngine.RegisterHook(engine.StatHook)
	srv := &http.Server{
		Addr:              ":8888",
		Handler:           oEngine,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}
	go func() {
		log.Println("listen", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("srv.ListenAndServe() err:", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
