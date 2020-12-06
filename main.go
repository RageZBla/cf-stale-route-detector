package main

import (
	"errors"
	"log"
	"os"

	"github.com/RageZBla/cf-stale-route-detector/cmd"
	"github.com/RageZBla/cf-stale-route-detector/commands"
)

var version = "unknown"

func main() {
	err := cmd.Main(os.Stdout, os.Stderr, version, os.Args)
	if err != nil {
		if errors.Is(err, commands.ErrStaleRouteDetected) {
			os.Exit(1)
		}
		log.Fatal(err)
	}
}
