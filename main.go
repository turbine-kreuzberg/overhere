package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
	bind = &cli.StringFlag{
		Name:    "bind-ip",
		Aliases: []string{"bind"},
		Usage:   "IP to bind.",
		Value:   "0.0.0.0",
		EnvVars: []string{"OVERHERE_PORT"},
	}
	port = &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Usage:   "Port to listen on.",
		Value:   53,
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
			bind,
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
	addr := c.String(bind.Name)
	port := c.Int(port.Name)
	verbose := c.Bool(verbose.Name)

	srv, err := overhere.NewServer(resolveTo, addr, port, verbose)
	if err != nil {
		return fmt.Errorf("setup dns server: %v", err)
	}

	log.Printf("listen on port: %s:%d", addr, port)

	err = srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("run dns server: %v", err)
	}

	ctx, cancel := context.WithTimeout(c.Context, 10*time.Second)
	defer cancel()

	err = srv.ShutdownContext(ctx)
	if err != nil {
		return fmt.Errorf("shutdown server: %v", err)
	}

	return nil
}
