package xsrf

import (
	"fmt"
	"html/template"
	"time"

	"github.com/Felamande/kiriadmin/modules/utils"
	"github.com/lunny/tango"
	"github.com/tango-contrib/cache"
)

//default values
const (
	CacheAdapter                  = "memory"
	Expiration      time.Duration = time.Minute * 30
	CacheGCInterval               = 120
	TokenKey                      = "_xsrf"
	RndGnrtorType                 = "sha1"
	DisposableKey                 = "once"
)

type Option struct {
	CacheAdapter    string
	RndGenerator    utils.RndGenerator
	Expiration      time.Duration
	CacheGCInterval int
	TokenKey        string
	RndGnrtorType   string
}

type XsrfHandler struct {
	Option
	cache        *cache.Caches
	rndGenerator utils.RndGenerator
}

func (x *XsrfHandler) Handle(ctx *tango.Context) {
	action := ctx.Action()
	checker, ok := action.(XsrfChecker)
	if !ok {
		ctx.Next()
		return
	}
	checker.SetHandler(x, ctx)

	switch ctx.Req().Method {
	case "POST":
		tokenFromReq := checker.GetTokenFromReq()
		if err := checker.ValidateXsrf(tokenFromReq); err != nil {

			checker.HandleXsrfErr(err)
			return
		}
		if checker.DisposableToken() {
			checker.HandleDisposable(tokenFromReq)
		}
	case "GET":
		token := checker.GenerateXsrfToken()
		if err := checker.PutToken(token); err != nil {
			checker.HandleXsrfErr(err)
			return
		}
	}

	ctx.Next()

}

func New(opts ...Option) *XsrfHandler {
	if len(opts) == 0 {
		return &XsrfHandler{
			Option: Option{
				CacheAdapter:  CacheAdapter,
				Expiration:    Expiration,
				TokenKey:      TokenKey,
				RndGnrtorType: RndGnrtorType,
			},
			rndGenerator: &utils.ShaGenerator{Type: RndGnrtorType},
			cache:        cache.New(cache.Options{Adapter: CacheAdapter, Interval: CacheGCInterval}),
		}
	}

	opt := opts[0]
	if opt.CacheAdapter == "" {
		opt.CacheAdapter = CacheAdapter
		// opt.cache =
	}
	if opt.CacheGCInterval == 0 {
		opt.CacheGCInterval = CacheGCInterval
	}
	if opt.Expiration == 0 {
		opt.Expiration = Expiration
	}
	if opt.RndGnrtorType == "" {
		opt.RndGenerator = &utils.ShaGenerator{Type: RndGnrtorType}
	}
	if opt.TokenKey == "" {
		opt.TokenKey = TokenKey
	}

	return &XsrfHandler{
		Option:       opt,
		cache:        cache.New(cache.Options{Adapter: opt.CacheAdapter, Interval: opt.CacheGCInterval}),
		rndGenerator: &utils.ShaGenerator{opt.RndGnrtorType},
	}
}

type XsrfChecker interface {
	SetHandler(xh *XsrfHandler, ctx *tango.Context)
	GenerateXsrfToken() string
	GetTokenFromReq() string
	ValidateXsrf(xsrfToken string) error
	PutToken(string) error
	Renew(token string) error
	DisposableToken() bool
	HandleDisposable(token string)
	ErrHandler
}

type Checker struct {
	xh    *XsrfHandler
	ctx   *tango.Context
	token string
}

var _ XsrfChecker = &Checker{}

func (c *Checker) CreateXsrfHTML() template.HTML {
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="%v" value="%v" />`,
		c.xh.TokenKey, c.token))
}

func (c *Checker) SetHandler(xh *XsrfHandler, ctx *tango.Context) {
	c.xh = xh
	c.ctx = ctx
}

func (c *Checker) GetTokenFromReq() string {
	c.ctx.Req().ParseForm()
	return c.ctx.Form(c.xh.TokenKey, "")
}

func (c *Checker) DisposableToken() bool {
	c.ctx.Req().ParseForm()
	v := c.ctx.Form(DisposableKey, "")
	return v == "1" || v == "true"
}

func (c *Checker) HandleDisposable(token string) {
	c.xh.cache.Delete(token)
}

func (c *Checker) GenerateXsrfToken() string {
	return c.xh.rndGenerator.GenerateRnd()
}

func (c *Checker) ValidateXsrf(token string) error {
	if token == "" {
		return XsrfError("empty token")

	}
	trueV := c.xh.cache.Get(token)
	if trueV == nil {
		return XsrfError("invalid token or expired")
	}
	return nil

}

func (c *Checker) PutToken(token string) error {
	c.token = token
	return c.xh.cache.Put(token, true, int64(c.xh.Expiration))
}

func (c *Checker) Renew(token string) error {
	trueV := c.xh.cache.Get(token)
	if trueV != nil {
		return nil
	}
	return c.xh.cache.Put(token, true, int64(c.xh.Expiration))

}

func (c *Checker) HandleXsrfErr(err error) {
	if xe, ok := err.(XsrfError); ok {
		c.ctx.Abort(401, xe.Error())
	} else {
		c.ctx.Abort(500, err.Error())
	}
}

type ErrHandler interface {
	HandleXsrfErr(err error)
}

type XsrfError string

func (e XsrfError) Error() string {
	return string(e)
}
