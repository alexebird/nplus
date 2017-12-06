package main

import (
	"os"

	npcli "github.com/alexebird/nplus/cli"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "nplus"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		npcli.TasksCliCommand(),
	}

	app.Run(os.Args)
}
