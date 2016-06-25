package settings

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Maxgis/tree"
	"github.com/kardianos/osext"
)

type staticCfg struct {
	RemoteDomain string `toml:"remote"`
	URLPrefix    string `toml:"urlprefix"`
	LocalRoot    string `toml:"local_root"`
	ListDir      bool   `toml:"listdir"`
}

type serverCfg struct {
	Port string `toml:"port"`
	Host string `toml:"host"`
}

type templateCfg struct {
	Home         string `toml:"home"`
	DelimesLeft  string `toml:"ldelime"`
	DelimesRight string `toml:"rdelime"`
	Charset      string `toml:"charset"`
	Reload       bool   `toml:"reload"`
}
type defaultVar struct {
	AppName string `toml:"appname"`
}

type adminCfg struct {
	Name string `toml:"name"`
	Enc  string `toml:"enc"`
	Rnd  string `toml:"rnd"`
}

type logCfg struct {
	Prefix string `toml:"prefix"`
	Path   string `toml:"path"`
	Format string `toml:"format"`
}

type timeCfg struct {
	ZoneString string `toml:"zone"`
}

type tlsCfg struct {
	Enable bool   `toml:"enable"`
	Cert   string `toml:"cert"`
	Key    string `toml:"key"`
}

type debugCfg struct {
	Port   int  `toml:"port"`
	Enable bool `toml:"enable"`
}

type captchaCfg struct {
	URLPrefix     string `toml:"urlprefix"`
	ChallengeNums int    `toml:"challenge_nums"`
	Width         int    `toml:"width"`
	Height        int    `toml:"height"`

	//expiration in seconds
	Expiration int64 `toml:"expiration"`
	Cache      struct {
		Adapter    string `toml:"adapter"`
		Config     string `toml:"config"`
		GCInterval int    `toml:"gc_interval"`
	} `toml:"cache"`
}

type xsrfCfg struct {
	RndGenerator string `toml:"rnd_generator"`
	//expiration in seconds
	Expiration int64 `toml:"expiration"`
	Cache      struct {
		Adapter    string `toml:"adapter"`
		Config     string `toml:"config"`
		GCInterval int    `toml:"gc_interval"`
	} `toml:"cache"`
}

type sessionCfg struct {
	//max age in hour
	MaxAge int64 `toml:"max_age"`
}

type setting struct {
	Static      staticCfg         `toml:"static"`
	Server      serverCfg         `toml:"server"`
	Template    templateCfg       `toml:"template"`
	DefaultVars defaultVar        `toml:"defaultvars"`
	Admin       adminCfg          `toml:"admin"`
	Log         logCfg            `toml:"log"`
	Time        timeCfg           `toml:"time"`
	TLS         tlsCfg            `toml:"tls"`
	Debug       debugCfg          `toml:"debug"`
	Headers     map[string]string `toml:"headers"`
	Captcha     captchaCfg        `toml:"captcha"`
	XSRF        xsrfCfg           `toml:"xsrf"`
	Session     sessionCfg        `toml:"session"`
	Params      map[string]string `toml:"params"`
}

var (
	Folder string
	// IsInit = false

	//GlobalSettings
	Static      staticCfg
	Server      serverCfg
	Template    templateCfg
	DefaultVars defaultVar
	Admin       adminCfg
	Log         logCfg
	TLS         tlsCfg
	Time        timeCfg
	Debug       debugCfg
	Captcha     captchaCfg
	XSRF        xsrfCfg
	Location    *time.Location
	Session     sessionCfg
	Params      map[string]string
	Headers     map[string]string
)

func init() {
	var err error
	Folder, err = osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}

}

var once sync.Once

func Init(cfgFile string) {
	once.Do(func() {
		settingStruct := setting{}
		b, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			panic(err)
		}
		if err := toml.Unmarshal(b, &settingStruct); err != nil {
			panic(err)
		}
		Location, err = time.LoadLocation(settingStruct.Time.ZoneString)
		if err != nil {
			fmt.Println(err)
			Location = time.UTC
		}

		Static = settingStruct.Static
		Server = settingStruct.Server

		Template = settingStruct.Template
		DefaultVars = settingStruct.DefaultVars
		Admin = settingStruct.Admin
		Log = settingStruct.Log
		Time = settingStruct.Time
		Headers = settingStruct.Headers
		Debug = settingStruct.Debug
		TLS = settingStruct.TLS
		Params = settingStruct.Params
		Captcha = settingStruct.Captcha
		XSRF = settingStruct.XSRF
		Session = settingStruct.Session

		if settingStruct.Debug.Enable {
			tree.Print(settingStruct)
		}

	})

}
