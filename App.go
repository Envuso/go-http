package main

import (
	"github.com/Envuso/go-ioc-container"
	"github.com/envuso/go-http/Http"
	"github.com/envuso/go-http/Routing"
	"github.com/envuso/go-http/lib/storage"
)

type App struct {
}

// var Container = Ioc.Container

func NewApp() *App {
	app := &App{}

	container.Bind(func() (storage.Storage, error) {
		storage.AddDiskConfig("test", storage.DiskConfig{
			Driver: "local",
			Root:   "Backup",
		})
		return storage.Disk("test")
	})

	container.Bind(func() SomeBullShitServiceContract {
		return NewSomeBullShitService("big yeet")
	})

	container.Bind(func() Routing.RouterContract {
		return Routing.Router
	})

	container.Bind(func(router Routing.RouterContract) Http.HttpContract {
		return Http.NewHttp().UsingRouter(router)
	})

	return app
}

var Application = NewApp()

func (app *App) Boot(addr string, implBoot interface{}) {

	if implBoot != nil {
		container.Call(implBoot)
	}

	var router Routing.RouterContract
	container.MakeTo(&router)

	router.Build()

	var server Http.HttpContract
	container.MakeTo(&server)

	err := server.Listen(addr)
	if err != nil {
		panic(err)
	}
}
