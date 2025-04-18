package main

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/canpacis/pacis/pages/pcpg/generator"
	pparser "github.com/canpacis/pacis/pages/pcpg/parser"
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
						return errors.New("a root directory is required for compiling")
					}

					root := ctx.Args().First()
					wd, err := os.Getwd()
					if err != nil {
						return err
					}
					absroot := path.Join(wd, root)

					list, err := pparser.ParseDir(path.Join(root, "app"))
					if err != nil {
						return err
					}
					assets, err := compile(root)
					if err != nil {
						return err
					}

					file, err := generator.CreateFile(list, assets)
					if err != nil {
						return err
					}

					content, err := generator.GenerateFile(file)
					if err != nil {
						return err
					}

					app := path.Join(absroot, "app")
					return os.WriteFile(path.Join(app, "app.gen.go"), content, 0o644)
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
