package main

import (
	"fmt"
	"os"
	_ "time/tzdata"

	"github.com/takumakei/runtil/app"
)

func main() {
	if err := app.App.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
