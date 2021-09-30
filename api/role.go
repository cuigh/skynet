package api

import (
	"time"

	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/store"
)

// RoleHandler encapsulates user related handlers.
type RoleHandler struct {
	Find   web.HandlerFunc `path:"/find" auth:"?" desc:"find role by id"`
	Search web.HandlerFunc `path:"/search" auth:"?" desc:"search roles"`
	Save   web.HandlerFunc `path:"/save" method:"post" auth:"role.save" desc:"create or update role"`
	Delete web.HandlerFunc `path:"/delete" method:"post" auth:"role.delete" desc:"delete role"`
}

// NewRole creates an instance of RoleHandler
func NewRole(s store.RoleStore) *RoleHandler {
	return &RoleHandler{
		Find:   roleFind(s),
		Search: roleSearch(s),
		Save:   roleSave(s),
		Delete: roleDelete(s),
	}
}

func roleFind(s store.RoleStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		role, err := s.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, role)
	}
}

func roleSave(s store.RoleStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		role := &store.Role{}
		err := ctx.Bind(role, true)
		if err == nil {
			if time.Time(role.CreateTime).IsZero() {
				err = s.Create(role)
			} else {
				err = s.Modify(role)
			}
		}
		return ajax(ctx, err)
	}
}

func roleSearch(s store.RoleStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		name := ctx.Query("name")
		roles, err := s.Search(name)
		if err != nil {
			return err
		}
		return success(ctx, roles)
	}
}

func roleDelete(s store.RoleStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		r := &store.Role{}
		err := ctx.Bind(r)
		if err == nil {
			err = s.Delete(r.ID)
		}
		return ajax(ctx, err)
	}
}
