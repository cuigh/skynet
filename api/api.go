package api

import (
	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/net/web"
)

type ajaxResult struct {
	Code int32       `json:"code"`
	Info string      `json:"info,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ajax(ctx web.Context, err error) error {
	if err != nil {
		return err
	}
	return ctx.JSON(ajaxResult{})
}

func success(ctx web.Context, data interface{}) error {
	return ctx.JSON(&ajaxResult{Data: data})
}

func init() {
	container.Put(NewSystem, container.Name("api.system"))
	container.Put(NewTask, container.Name("api.task"))
	container.Put(NewJob, container.Name("api.job"))
	container.Put(NewUser, container.Name("api.user"))
	container.Put(NewRole, container.Name("api.role"))
	container.Put(NewConfig, container.Name("api.config"))
}
