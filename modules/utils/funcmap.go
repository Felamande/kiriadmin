package utils

import (
	"fmt"
	"html/template"
	"path"
	"path/filepath"
	"strings"

	"github.com/Felamande/kiriadmin/settings"
)

var AssetTpl = map[string]string{
	"css": `<link rel="stylesheet" href="%s" type="text/css" />`,
	"js":  `<script src="%s"></script>`,
}

func AssetLocal(typ, src string) template.HTML {

	return template.HTML(fmt.Sprintf(AssetTpl[typ], path.Join(settings.Static.URLPrefix, typ, src)))
	// return path.Join(settings.Static.VirtualRoot, "js", src)
}

func AssetRemote(typ, src string) template.HTML {
	return template.HTML(fmt.Sprintf(AssetTpl[typ], "https://"+path.Join(settings.Static.RemoteDomain, typ, src)))
}

func DefaultFuncs() template.FuncMap {
	// s, err := compress.LoadJsonConf(Abs(settings.Static.CompressDef), true, settings.Server.Host)
	// if err != nil {
	// 	panic(err)

	// }
	return template.FuncMap{
		"AssetLocal":  AssetLocal,
		"AssetRemote": AssetRemote,
		// "CompressCss": s.Css.CompressCss,
		// "CompressJs":  s.Js.CompressJs,
	}
}

func Ext(filename string) string {
	clean := strings.Split(filename, "?")[0]
	return filepath.Ext(clean)
}
