package api

import (
	"github.com/cuigh/auxo/app/ioc"
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
	ioc.Put(NewSystem, ioc.Name("api.system"))
	ioc.Put(NewTask, ioc.Name("api.task"))
	ioc.Put(NewJob, ioc.Name("api.job"))
	ioc.Put(NewUser, ioc.Name("api.user"))
	ioc.Put(NewRole, ioc.Name("api.role"))
	ioc.Put(NewConfig, ioc.Name("api.config"))
}
