package utils

import (
	"fmt"
	"html/template"
	"path"

	"github.com/Felamande/kiriadmin/settings"
)

func AssetLocal(typ, src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, path.Join(settings.Static.VirtualRoot, typ, src)))
	// return path.Join(settings.Static.VirtualRoot, "js", src)
}

func AssetRemote(typ, src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<link rel="stylesheet" href="%s" type="text/css" />`, "https://"+path.Join(settings.Static.RemoteRoot, "css", src)))
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
