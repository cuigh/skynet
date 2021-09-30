package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/app/flag"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/data/valid"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/net/web/filter"
	_ "github.com/cuigh/skynet/api"
	_ "github.com/cuigh/skynet/lock"
	"github.com/cuigh/skynet/runner"
	"github.com/cuigh/skynet/schedule"
	_ "github.com/cuigh/skynet/store"
)

var (
	//go:embed ui/dist
	webFS embed.FS
)

func main() {
	app.Name = "Skynet"
	app.Version = "0.2"
	app.Desc = "A distributed task scheduler"
	app.Flags.Register(flag.All)
	app.Action = entry
	app.Start()
}

func entry(*app.Context) error {
	// 启动调度器
	app.Ensure(container.Call(func(s *schedule.Scheduler) {
		go s.Start()
	}))

	// 启动网站
	app.Run(createWebServer())
	return nil
}

func createWebServer() *web.Server {
	ws := web.Auto()
	ws.Validator = &valid.Validator{}
	ws.ErrorHandler.Default = handleError
	ws.Use(filter.NewRecover())
	ws.Static("/", http.FS(loadWebFS()), "index.html")

	g := ws.Group("/api", findFilters("authenticator", "authorizer")...)
	g.Handle("/system", container.Find("api.system"))
	g.Handle("/task", container.Find("api.task"))
	g.Handle("/job", container.Find("api.job"))
	g.Handle("/user", container.Find("api.user"))
	g.Handle("/role", container.Find("api.role"))
	g.Handle("/config", container.Find("api.config"))

	// runner testing
	ws.Post("/task/execute", runner.HandleExecute, web.WithAuthorize(web.AuthAnonymous))
	ws.Post("/task/split", runner.HandleSplit, web.WithAuthorize(web.AuthAnonymous))

	return ws
}

func loadWebFS() fs.FS {
	sub, err := fs.Sub(webFS, "ui/dist")
	if err != nil {
		panic(err)
	}
	return sub
}

func handleError(ctx web.Context, err error) {
	var (
		status       = http.StatusInternalServerError
		code   int32 = 1
	)

	if e, ok := err.(*web.Error); ok {
		status = e.Status()
	}
	if e, ok := err.(*errors.CodedError); ok {
		code = e.Code
	}

	err = ctx.Status(status).JSON(data.Map{"code": code, "info": err.Error()})
	if err != nil {
		ctx.Logger().Error(err)
	}
}

func findFilters(names ...string) []web.Filter {
	var filters []web.Filter
	for _, name := range names {
		filters = append(filters, container.Find(name).(web.Filter))
	}
	return filters
}