package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-vcf-tools"
	app.Version = Version
	app.Usage = "Simple VCF tools written in golang"
	app.Author = "knmkr"
	app.Email = "knmkr3gma@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
