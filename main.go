package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/threeaccents/digolang/repl"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the digo programming language.\n Feel free to type in commands\n", u.Username)
	repl.Start(os.Stdin, os.Stdout)
}
