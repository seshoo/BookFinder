package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/seshoo/bookFinder/internal/app"
	"os"
)

var revision = "0.0.1"

var exitFunc = os.Exit

func main() {
	fmt.Printf("Book finder: %s\n", revision)

	var opts app.Options
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	if _, err := p.Parse(); err != nil {
		exitFunc(1)
	}

	if err := app.Run(opts); err != nil {
		fmt.Printf("[WARN] %v", err)
	}
}
