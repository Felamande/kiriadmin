//github.com/go-xorm/dbweb/middlewares/auth.go

package auth

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
)

const (
	adminTokenKey = "adminToken"
)

type auther interface {
	AskAuth() bool
	SetToken(string)
	SetSession(*session.Session)
}

type Auther struct {
	token string
	s     *session.Session
}

func (a *Auther) Token() string {
	return a.token
}

func (a *Auther) Logout() {
	a.s.Del(adminTokenKey)
	a.s.Release()
}

func (a *Auther) AskAuth() bool {
	return true
}

func (a *Auther) LoginWithToken(token string) {
	a.SetToken(token)
	a.s.Set(adminTokenKey, token)
}

func (a *Auther) SetToken(token string) {
	a.token = token
}

func (a *Auther) SetSession(sess *session.Session) {
	a.s = sess
}

func Auth(redirct string, sessions *session.Sessions) tango.HandlerFunc {
	return func(ctx *tango.Context) {
		if auther, ok := ctx.Action().(auther); ok {
			s := sessions.Session(ctx.Req(), ctx.ResponseWriter)
			auther.SetSession(s)
			if adminToken := s.Get(adminTokenKey); adminToken != nil {
				auther.SetToken(adminToken.(string))
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
