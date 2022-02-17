package main

import (
	"fmt"
	"log"
	"os"

	overhere "github.com/turbine-kreuzberg/overhere/pkg"
	"github.com/urfave/cli/v2"
)

var (
	resolveTo = &cli.StringFlag{
		Name:    "resolve-to",
		Aliases: []string{"ip"},
		Usage:   "IP to resolve to. Autodetected by default.",
		EnvVars: []string{"OVERHERE_RESOLVE_TO"},
	}
	port = &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Usage:   "Port to listen on.",
		Value:   15353,
		EnvVars: []string{"OVERHERE_PORT"},
	}
	verbose = &cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		Usage:   "Verbosity of logging.",
		EnvVars: []string{"OVERHERE_VERBOSE"},
	}
	app = &cli.App{
		Name:  "Overhere",
		Usage: "A very minimal DNS server for development purposes.",
		Flags: []cli.Flag{
			resolveTo,
			port,
			verbose,
		},
		Action: runServer,
	}
)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(c *cli.Context) error {
	resolveTo := c.String(resolveTo.Name)
	port := c.Int(port.Name)
	verbose := c.Bool(verbose.Name)

	srv, err := overhere.NewServer(resolveTo, port, verbose)
	if err != nil {
		return fmt.Errorf("setup dns server: %v", err)
	}

	log.Printf("listen on port: %d", port)

	err = srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("run dns server: %v", err)
	}

	srv.Shutdown()

	return nil
}
