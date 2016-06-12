package main

import (
	"os"
	"time"

	"github.com/Felamande/kiriadmin/settings"
	"github.com/Felamande/lotdb/routers/toolate"
	"github.com/lunny/tango"

	"github.com/Felamande/kiriadmin/models"

	//modules
	"github.com/Felamande/kiriadmin/modules/log"
	"github.com/Felamande/kiriadmin/modules/utils"

	//middlewares
	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/middlewares/header"
	timemw "github.com/Felamande/kiriadmin/middlewares/time"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/cache"
	"github.com/tango-contrib/captcha"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"
	"github.com/tango-contrib/xsrf"

	//routers
	"github.com/Felamande/kiriadmin/routers/debug"
	"github.com/Felamande/kiriadmin/routers/home"
	"github.com/Felamande/kiriadmin/routers/login"
)

func init() {
	settings.Init("./settings/settings.toml")
	models.Init()
}

func main() {
	l := log.New(os.Stdout, settings.Log.Prefix, log.Llevel|log.LstdFlags)
	l.SetLocation(settings.Time.Location)
	t := tango.NewWithLog(l)

	sess := session.New(session.Options{
		MaxAge: time.Hour * 24,
	})

	CaptchaCache := cache.New(cache.Options{
		Adapter:  "memory",
		Interval: 120,
	})

	t.Use(
		new(timemw.TimeHandler),
		binding.Bind(),
		tango.Recovery(false),
		tango.Compresses([]string{}),
		tango.Return(),
		tango.Param(),
		tango.Contexts(),
		renders.New(renders.Options{
			Reload:      settings.Template.Reload,
			Directory:   settings.Template.Home,
			Charset:     settings.Template.Charset,
			DelimsLeft:  settings.Template.DelimesLeft,
			DelimsRight: settings.Template.DelimesRight,
			Funcs:       utils.DefaultFuncs(),
		}),
		events.Events(),
		header.CustomHeaders(),
		captcha.New(captcha.Options{
			URLPrefix:        "/captcha/",  // URL prefix of getting captcha pictures.
			FieldIdName:      "captcha_id", // Hidden input element ID.
			FieldCaptchaName: "captcha",    // User input value element name in request form.
			ChallengeNums:    6,            // Challenge number.
			Width:            240,          // Captcha image width.
			Height:           80,           // Captcha image height.
			Expiration:       600,          // Captcha expiration time in seconds.
			CachePrefix:      "captcha_",   // Cache key prefix captcha characters.
			Caches:           CaptchaCache,
		}),
		auth.Auth("/login", sess),
		xsrf.New(time.Minute),
	)

	t.Get("/", new(home.HomeRouter))
	t.Any("/login", new(login.LoginRouter))
	t.Post(toolate.Url, new(toolate.TooLateRouter))

	if settings.Debug.Enable {
		go debug.On(settings.Debug.Port)
		l.Infof("enable debug at port %d\n", settings.Debug.Port)
	}

	if settings.TLS.Enable {
		t.RunTLS(settings.TLS.Cert, settings.TLS.Key, settings.Server.Host)
	} else {
		t.Run(settings.Server.Host)
	}

}
