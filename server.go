package main

import (
	"os"
	"time"

	"github.com/Felamande/kiriadmin/settings"
	"github.com/lunny/tango"

	//models
	"github.com/Felamande/kiriadmin/models"

	//modules
	"github.com/Felamande/kiriadmin/modules/log"
	"github.com/Felamande/kiriadmin/modules/utils"

	//middlewares
	"github.com/Felamande/kiriadmin/middlewares/auth"
	"github.com/Felamande/kiriadmin/middlewares/header"
	timemw "github.com/Felamande/kiriadmin/middlewares/time"
	"github.com/Felamande/kiriadmin/middlewares/xsrf"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/cache"
	_ "github.com/tango-contrib/cache-nodb"
	"github.com/tango-contrib/captcha"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
	"github.com/tango-contrib/session"

	//routers
	"github.com/Felamande/kiriadmin/routers/debug"
	"github.com/Felamande/kiriadmin/routers/editor"
	"github.com/Felamande/kiriadmin/routers/home"
	"github.com/Felamande/kiriadmin/routers/login"
)

func init() {
	settings.Init("./settings/settings.toml")
	models.Init()
}

func main() {
	l := log.New(os.Stdout, settings.Log.Prefix, log.Llevel|log.LstdFlags)
	l.SetLocation(settings.Location)
	tgo := tango.NewWithLog(l)

	sess := session.New(session.Options{
		MaxAge: time.Duration(settings.Session.MaxAge * int64(time.Hour)),
	})

	CaptchaCache := cache.New(cache.Options{
		Adapter:       settings.Captcha.Cache.Adapter,
		Interval:      settings.Captcha.Cache.GCInterval,
		AdapterConfig: settings.Captcha.Cache.Config,
	})

	XSRFCache := cache.New(cache.Options{
		Adapter:       settings.XSRF.Cache.Adapter,
		Interval:      settings.XSRF.Cache.GCInterval,
		AdapterConfig: settings.XSRF.Cache.Config,
	})

	tgo.Use(
		new(timemw.TimeHandler),
		tango.Static(tango.StaticOptions{
			RootPath: settings.Static.LocalRoot,
			Prefix:   settings.Static.URLPrefix,
			ListDir:  settings.Static.ListDir,
		}),
		tango.Recovery(false),
		tango.Return(),
		tango.Param(),
		tango.Contexts(),
		binding.Bind(),
		events.Events(),

		renders.New(renders.Options{
			Reload:      settings.Template.Reload,
			Directory:   settings.Template.Home,
			Charset:     settings.Template.Charset,
			DelimsLeft:  settings.Template.DelimesLeft,
			DelimsRight: settings.Template.DelimesRight,
			Funcs:       utils.DefaultFuncs(),
		}),
		captcha.New(captcha.Options{
			URLPrefix:        settings.Captcha.URLPrefix,     // URL prefix of getting captcha pictures.
			FieldIdName:      "captcha_id",                   // Hidden input element ID.
			FieldCaptchaName: "captcha",                      // User input value element name in request form.
			ChallengeNums:    settings.Captcha.ChallengeNums, // Challenge number.
			Width:            settings.Captcha.Width,         // Captcha image width.
			Height:           settings.Captcha.Height,        // Captcha image height.
			Expiration:       settings.Captcha.Expiration,    // Captcha expiration time in seconds.
			CachePrefix:      "captcha_",                     // Cache key prefix captcha characters.
			Caches:           CaptchaCache,
		}),
		auth.Auth("/login", sess),
		// xsrf.New(time.Minute),
		xsrf.New(xsrf.Option{
			Cache:        XSRFCache,
			RndGenerator: &utils.ShaGenerator{Type: settings.XSRF.RndGenerator},
			Expiration:   settings.XSRF.Expiration,
		}),
		header.CustomHeaders(),
	)

	tgo.Get("/", new(home.HomeRouter))
	tgo.Any("/login", new(login.LoginRouter))
	tgo.Get("/logout", new(login.LogoutRouter))
	tgo.Group("/editor", func(g *tango.Group) {
		g.Get("/", new(editor.EditorHome))
		g.Post("/preview", new(editor.PreviewRouter))
	})

	if settings.Debug.Enable {
		go debug.On(settings.Debug.Port)
		l.Infof("enable debug at port %d\n", settings.Debug.Port)
	}

	if settings.TLS.Enable {
		tgo.RunTLS(settings.TLS.Cert, settings.TLS.Key, settings.Server.Host)
	} else {
		tgo.Run(settings.Server.Host)
	}

}
