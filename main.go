package main

import (
	"flag"
	"github.com/sagernet/sagernetdesktop/app"
	"os"
)

//go:generate goversioninfo --platform-specific

func main() {
	fs := flag.NewFlagSet("SagerNet", flag.ExitOnError)
	headless := fs.Bool("h", false, "headless mode")
	_ = fs.Parse(os.Args[1:])

	if *headless {
		println("Hello World :)")
	} else {
		app.Launch()
	}
}
