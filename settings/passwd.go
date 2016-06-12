//+build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/Felamande/kiriadmin/modules/utils"
	"github.com/peterh/liner"
)

var encType = flag.String("type", "sha256", "encrypt type.")
var rndlen = flag.Int64("rndlen", 16, "random string length")

func main() {
	flag.Parse()
	line := liner.NewLiner()
	pwd1, err := line.PasswordPrompt("password:")
	if err != nil {
		fmt.Println(err)
		return
	}
	pwd2, err := line.PasswordPrompt("password:")
	if err != nil {
		fmt.Println(err)
		return
	}
	if pwd1 != pwd2 {
		fmt.Println("passwords don't agree.")
		return
	}

	rands := utils.RandomString(*rndlen)
	encFunc, exist := utils.Encrypt[*encType]
	if !exist {
		fmt.Println("unsupported ecrypt type:", *encType)
	}
	enc := encFunc(rands + pwd1)
	fmt.Printf(`enc ="%s"
`, enc)
	fmt.Printf(`rnd ="%s"
`, rands)

}
