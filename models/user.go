package models

import (
	"github.com/Felamande/kiriadmin/modules/utils"
	"github.com/Felamande/kiriadmin/settings"
)

func Init() {
	Admin = &User{
		Name: settings.Admin.Name,
		Enc:  settings.Admin.Enc,
		Rnd:  settings.Admin.Rnd,
	}
}

type User struct {
	Name string
	Enc  string
	Rnd  string
}

type UserForm struct {
	Name string `form:"name"`
	Pwd  string `form:"pwd"`
}

func GetUserByName(name string) *User {
	if name != Admin.Name {
		return nil
	}
	return Admin
}

func (u *User) Auth(pwd string) bool {
	return AuthUser(u, pwd)
}

func AuthUser(u *User, pwd string) bool {
	return u.Enc == utils.Encrypt["sha256"](u.Rnd+pwd)
}

var Admin *User
