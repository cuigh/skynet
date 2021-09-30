package api

import (
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/store"
)

// ConfigHandler encapsulates config related handlers.
type ConfigHandler struct {
	Find web.HandlerFunc `path:"/find" auth:"?" desc:"find config by id"`
	Save web.HandlerFunc `path:"/save" method:"post" auth:"config.edit" desc:"modify config"`
}

// NewConfig creates an instance of ConfigHandler
func NewConfig(store store.ConfigStore) *ConfigHandler {
	return &ConfigHandler{
		Find: configFind(store),
		Save: configSave(store),
	}
}

func configFind(s store.ConfigStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		options, err := s.Find(id)
		if err != nil {
			return err
		}

		m := make(map[string]string)
		for _, opt := range options {
			m[opt.Name] = opt.Value
		}
		return success(ctx, m)
	}
}

func configSave(s store.ConfigStore) web.HandlerFunc {
	type Args struct {
		Id      string            `json:"id"`
		Options map[string]string `json:"options"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		var options data.Options
		for k, v := range args.Options {
			options = append(options, data.Option{Name: k, Value: v})
		}

		err = s.Save(args.Id, options)
		return ajax(ctx, err)
	}
}
