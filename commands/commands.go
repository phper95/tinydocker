package commands

import (
	"context"
	"github.com/phper95/tinydocker/container"
	"github.com/urfave/cli/v3"
	"log"
)

var RunCommand = &cli.Command{
	Name:  "run",
	Usage: "Run a command in a new container",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "it",
			Usage: "Interactive mode with pseudo-TTY",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) == 0 {
			log.Fatal("No command specified")
		}
		return container.Run(args, cmd.Bool("it"))
	},
}
