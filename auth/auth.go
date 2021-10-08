package auth

import (
	"net/http"
	"time"

	"github.com/cuigh/auxo/app/container"
	"github.com/cuigh/auxo/cache"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/skynet/store"
)

type Authorizer struct {
	us store.UserStore
	rv *cache.Value
}

func NewAuthorizer(us store.UserStore, rs store.RoleStore) web.Filter {
	v := cache.Value{
		TTL: 5 * time.Minute,
		Load: func() (interface{}, error) {
			return loadRolePerms(rs)
		},
	}
	return &Authorizer{us: us, rv: &v}
}

// Apply implements `web.Filter` interface.
func (a *Authorizer) Apply(next web.HandlerFunc) web.HandlerFunc {
	return func(ctx web.Context) error {
		auth := ctx.Handler().Authorize()

		// allow anonymous
		if auth == "" || auth == web.AuthAnonymous {
			return next(ctx)
		}

		user := ctx.User()
		if user == nil || user.Anonymous() {
			return web.NewError(http.StatusUnauthorized, "You are not logged in")
		}

		if auth != web.AuthAuthenticated && !a.check(user, auth) {
			return web.NewError(http.StatusForbidden, "You do not have access to this resource")
		}
		return next(ctx)
	}
}

func (a *Authorizer) check(user web.User, auth string) bool {
	u, err := a.us.Find(user.ID())
	if err != nil {
		log.Get("auth").Error(err)
		return false
	}

	if u.Admin {
		return true
	}

	rp := a.rv.MustGet(true).(RolePerms)
	for _, role := range u.Roles {
		if rp[role].Contains(auth) {
			return true
		}
	}
	return false
}

type RolePerms map[string]data.Map

func loadRolePerms(rs store.RoleStore) (RolePerms, error) {
	roles, err := rs.Search("")
	if err != nil {
		return nil, err
	}

	rp := make(RolePerms)
	for _, role := range roles {
		m := make(data.Map)
		for _, p := range role.Perms {
			m.Set(p, nil)
		}
		rp[role.ID] = m
	}
	return rp, nil
}

func init() {
	container.Put(NewAuthorizer, container.Name("authorizer"))
	container.Put(NewAuthenticator, container.Name("authenticator"))
}
