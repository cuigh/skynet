package api

import (
	"github.com/cuigh/auxo/app/ioc"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security/passwd"
	"github.com/cuigh/skynet/auth"
	"github.com/cuigh/skynet/store"
	"strings"
)

// UserHandler encapsulates user related handlers.
type UserHandler struct {
	SignIn         web.HandlerFunc `path:"/sign-in" method:"post" auth:"*" desc:"user sign in"`
	Search         web.HandlerFunc `path:"/search" auth:"?" desc:"search users"`
	Find           web.HandlerFunc `path:"/find" auth:"?" desc:"find user by id"`
	Fetch          web.HandlerFunc `path:"/fetch" auth:"?" desc:"fetch users by ids"`
	ModifyPassword web.HandlerFunc `path:"/modify-password" method:"post" auth:"?" desc:"modify password"`
	ModifyProfile  web.HandlerFunc `path:"/modify-profile" method:"post" auth:"?" desc:"modify profile"`
	Save           web.HandlerFunc `path:"/save" method:"post" auth:"user.edit" desc:"create or update user"`
	SetStatus      web.HandlerFunc `path:"/set-status" method:"post" auth:"user.edit" desc:"set user status"`
	//SignUp         web.HandlerFunc `path:"/sign-up" method:"post" auth:"*" desc:"user sign up"`
	//RefreshToken   web.HandlerFunc `method:"post" path:"/refresh-token" auth:"?" desc:"refresh user token"`
}

// NewUser creates an instance of UserHandler
func NewUser(s store.UserStore) *UserHandler {
	return &UserHandler{
		SignIn:         userSignIn(s),
		Search:         userSearch(s),
		Find:           userFind(s),
		Fetch:          userFetch(s),
		ModifyPassword: userModifyPassword(s),
		ModifyProfile:  userModifyProfile(s),
		Save:           userSave(s),
		SetStatus:      userSetStatus(s),
		//SignUp:         userSignUp,
		//RefreshToken:   userRefreshToken,
	}
}

func userSignIn(store store.UserStore) web.HandlerFunc {
	type SignInArgs struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	return func(ctx web.Context) error {
		args := &SignInArgs{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		user, err := store.FindByName(args.Name)
		if err != nil {
			return err
		} else if user != nil && passwd.Validate(args.Password, user.Password, user.Salt) {
			jwt := ioc.Find[*auth.JWT]("authenticator")
			token, err := jwt.CreateToken(user.Id, user.Name)
			if err != nil {
				return err
			}
			return success(ctx, data.Map{
				"token": token,
				"id":    user.Id,
				"name":  user.Name,
			})
		}
		return errors.New("user not exists or password incorrect")
	}
}

//func userSignUp(c web.Context) error {
//	return errors.NotImplemented
//}

func userSave(s store.UserStore) web.HandlerFunc {
	const initialPassword = "754321"

	return func(ctx web.Context) error {
		u := &store.User{}
		err := ctx.Bind(u, true)
		if err == nil {
			if u.Id == "" {
				u.Password, u.Salt, err = passwd.Generate(initialPassword)
				if err == nil {
					err = s.Create(u)
				}
			} else {
				err = s.Modify(u)
			}
		}
		return ajax(ctx, err)
	}
}

//func userRefreshToken(ctx web.Context) error {
//	jwt := ioc.Find[*auth.JWT]("authenticator")
//	token, err := jwt.CreateToken(ctx.User().ID(), ctx.User().Name())
//	if err != nil {
//		return err
//	}
//	return success(ctx, data.Map{
//		"token": token,
//		"id":    ctx.User().ID(),
//		"name":  ctx.User().Name(),
//	})
//}

func userSearch(s store.UserStore) web.HandlerFunc {
	type Args struct {
		Name      string `json:"name"`
		LoginName string `json:"login_name"`
		Filter    string `json:"filter"`
		PageIndex int64  `json:"page_index"`
		PageSize  int64  `json:"page_size"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		users, total, err := s.Search(args.Name, args.LoginName, args.Filter, args.PageIndex, args.PageSize)
		if err != nil {
			return err
		}
		return success(ctx, data.Map{"items": users, "total": total})
	}
}

func userFind(s store.UserStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		id := ctx.Query("id")
		if id == "" {
			id = ctx.User().ID()
		}
		user, err := s.Find(id)
		if err != nil {
			return err
		}
		return success(ctx, user)
	}
}

func userFetch(s store.UserStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		ids := ctx.Query("ids")
		if ids == "" {
			return errors.New("missing argument 'ids'")
		}
		users, err := s.Fetch(strings.Split(ids, ","))
		if err != nil {
			return err
		}
		return success(ctx, users)
	}
}

func userSetStatus(s store.UserStore) web.HandlerFunc {
	type Args struct {
		ID     string `json:"id"`
		Status int32  `json:"status"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err == nil {
			err = s.SetStatus(args.ID, args.Status)
		}
		return ajax(ctx, err)
	}
}

func userModifyPassword(s store.UserStore) web.HandlerFunc {
	type Args struct {
		OldPassword string `json:"old_pwd"`
		NewPassword string `json:"new_pwd"`
	}

	return func(ctx web.Context) error {
		args := &Args{}
		err := ctx.Bind(args)
		if err != nil {
			return err
		}

		id := ctx.User().ID()
		user, err := s.Find(ctx.User().ID())
		if err != nil {
			return err
		} else if user == nil {
			return errors.Format("未找到用户: %d", id)
		}

		if !passwd.Validate(args.OldPassword, user.Password, user.Salt) {
			return errors.New("密码不正确")
		}

		pwd, salt, err := passwd.Generate(args.NewPassword)
		if err == nil {
			err = s.SetPassword(id, pwd, salt)
		}
		return ajax(ctx, err)
	}
}

func userModifyProfile(s store.UserStore) web.HandlerFunc {
	return func(ctx web.Context) error {
		u := &store.User{}
		err := ctx.Bind(u, true)
		if err == nil {
			err = s.ModifyProfile(u)
		}
		return ajax(ctx, err)
	}
}
