package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/log"
	"github.com/cuigh/auxo/net/web"
	"github.com/cuigh/auxo/security"
	"github.com/cuigh/auxo/util/cast"
	"github.com/dgrijalva/jwt-go"
)

var ErrNoNeedRefresh = errors.New("no need to refresh")

type JWT struct {
	Schema      string
	Sources     data.Options
	KeyFunc     jwt.Keyfunc
	Identifier  func(token *jwt.Token) web.User
	tokenExpiry int64
}

func NewAuthenticator() web.Filter {
	key := config.GetString("token_key")
	expiry := config.GetDuration("token_expiry")
	if key == "" {
		key = "skynet"
	}
	if expiry == 0 {
		expiry = 30 * time.Minute
	}

	return &JWT{
		tokenExpiry: int64(expiry.Seconds()),
		Schema:      "Bearer",
		Sources: data.Options{
			{Name: "header", Value: web.HeaderAuthorization},
		},
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			// TODO: use user salt as key
			return []byte(key), nil
		},
		Identifier: func(token *jwt.Token) web.User {
			claims := token.Claims.(jwt.MapClaims)
			return security.NewUser(cast.ToString(claims["sub"]), cast.ToString(claims["name"]))
		},
	}
}

func (j *JWT) Apply(next web.HandlerFunc) web.HandlerFunc {
	if j.KeyFunc == nil {
		panic("KeyFunc is required")
	}
	if j.Schema == "" {
		j.Schema = "Bearer"
	}
	if len(j.Sources) == 0 {
		j.Sources = data.Options{
			{Name: "header", Value: web.HeaderAuthorization},
		}
	}
	if j.Identifier == nil {
		j.Identifier = func(token *jwt.Token) web.User {
			claims := token.Claims.(jwt.MapClaims)
			return security.NewUser(cast.ToString(claims["sub"]), cast.ToString(claims["name"]))
		}
	}

	logger := log.Get("auth")
	return func(ctx web.Context) error {
		ts := j.extractToken(ctx)
		if ts != "" {
			token, err := jwt.Parse(ts, j.KeyFunc)
			if err != nil {
				logger.Debugf("failed to parse token: %s", err)
			} else {
				user := j.Identifier(token)
				ctx.SetUser(user)
				if ts, err := j.refreshToken(user, token); err == nil {
					ctx.SetHeader(web.HeaderAuthorization, ts)
				}
			}
		}

		return next(ctx)
	}
}

func (j *JWT) extractToken(ctx web.Context) (token string) {
	for _, src := range j.Sources {
		switch src.Name {
		case "header":
			token = ctx.Header(src.Value)
			if strings.HasPrefix(token, j.Schema) {
				return token[len(j.Schema)+1:]
			}
		case "cookie":
			if cookie, err := ctx.Cookie(src.Value); err == nil {
				token = cookie.Value
			}
		case "form":
			token = ctx.Form(src.Value)
		case "query":
			token = ctx.Query(src.Value)
		}
		if token != "" {
			return
		}
	}
	return
}

func (j *JWT) refreshToken(user web.User, token *jwt.Token) (string, error) {
	claims := token.Claims.(jwt.MapClaims)
	expiry := cast.ToInt64(claims["exp"])
	now := time.Now().Unix()
	// refresh token when remaining expiry is less than 5 minutes
	if (expiry - now) < 5*60 {
		ts, err := j.CreateToken(user.ID(), user.Name())
		if err != nil {
			return "", err
		}
		return ts, nil
	}
	return "", ErrNoNeedRefresh
}

func (j *JWT) CreateToken(id, name string) (string, error) {
	// https://www.iana.org/assignments/jwt/jwt.xhtml
	//iss: jwt签发者
	//sub: jwt所面向的用户
	//aud: 接收jwt的一方
	//exp: jwt的过期时间，这个过期时间必须要大于签发时间
	//nbf: 定义在什么时间之前，该jwt都是不可用的.
	//iat: jwt的签发时间
	//jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"name": name,
		"sub":  id,
		"iat":  now,
		"exp":  now + j.tokenExpiry,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := j.KeyFunc(token)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}
