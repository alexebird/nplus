package cli

import (
	"log"

	nomad "github.com/alexebird/nplus/nomad"
	//"github.com/davecgh/go-spew/spew"
	"github.com/urfave/cli"
)

func tasksCliAction(c *cli.Context) {
	client, err := nomad.Client()
	if err != nil {
		log.Fatal(err)
	}

	allocs, err := client.Allocations()
	if err != nil {
		log.Fatal(err)
	}
	//spew.Dump(allocs)

	if c.Bool("all") {
		nomad.PrintTasksTableLong(allocs)
	} else {
		nomad.PrintTasksTableShort(allocs)
	}
}

func TasksCliCommand() cli.Command {
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "dont do default filtering",
		},
		cli.StringFlag{
			Name:  "filter, f",
			Value: "",
			Usage: "filter by instance name",
		},
	}

	return cli.Command{
		Name:   "tasks",
		Action: tasksCliAction,
		Flags:  flags,
		UseShortOptionHandling: true,
	}
}
