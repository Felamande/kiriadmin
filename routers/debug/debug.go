package debug

import (
	"net/http"
	_ "net/http/pprof"
	"strconv"
)

func On(port int) {
	http.ListenAndServe("127.0.0.1:"+strconv.Itoa(port), nil)
}
