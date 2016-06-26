//github.com/go-xorm/dbweb/middlewares/auth.go

package auth

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
)

type Auther interface {
	AskAuth() bool
	SetSession(*session.Session)
	TokenKey() string

	setToken(interface{})
}

type auther struct {
	s *session.Session
}

func (a *auther) SetSession(sess *session.Session) {
	a.s = sess
}

//check implimentation
var _ Auther = &AuthUser{}

type AuthUser struct {
	token interface{}
	auther
}

func (a *AuthUser) Token() interface{} {
	return a.token
}

func (a *AuthUser) TokenKey() string {
	return "UserTokenKey"
}

func (a *AuthUser) Logout() {
	a.s.Del(a.TokenKey())
	a.s.Release()
}

func (a *AuthUser) AskAuth() bool {
	return true
}

func (a *AuthUser) LoginWithToken(token interface{}) {
	a.setToken(token)
	a.s.Set(a.TokenKey(), token)
}

func (a *AuthUser) setToken(token interface{}) {
	a.token = token
}

//check implimentation
var _ Auther = &AuthAdmin{}

type AuthAdmin struct {
	AuthUser
}

func (a *AuthAdmin) TokenKey() string {
	return "adminKey"
}

func (a *AuthAdmin) LoginWithToken(token interface{}) {
	a.setToken(token)
	a.s.Set(a.TokenKey(), token)
}

func Auth(redirct string, sessions *session.Sessions) tango.HandlerFunc {
	return func(ctx *tango.Context) {

		if auther, ok := ctx.Action().(Auther); ok {
			s := sessions.Session(ctx.Req(), ctx.ResponseWriter)
			auther.SetSession(s)
			if token := s.Get(auther.TokenKey()); token != nil {
				auther.setToken(token)
			} else {
				if auther.AskAuth() {
					ctx.Redirect(redirct)
					return
				}
			}
		}
		ctx.Next()
	}
}
