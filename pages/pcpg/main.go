package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pcpg",
		Usage: "pacis pages cli for bundling your assets and managing binaries for tooling",
		Commands: []*cli.Command{
			{
				Name:    "compile",
				Aliases: []string{"c"},
				Usage:   "bundle your assets and compile your go code to create a router",
				Action: func(ctx *cli.Context) error {
					if !ctx.Args().Present() {
						return errors.New("a target directory is required for compiling")
					}

					target := ctx.Args().First()
					if err := compile(target); err != nil {
						return err
					}
					gen, err := scan(target)
					if err != nil {
						return err
					}
					return generate(target, gen)
				},
			},
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "install/update dependencies for pacis pages cli",
				Action: func(ctx *cli.Context) error {
					return install()
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
