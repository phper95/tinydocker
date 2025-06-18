package main

import (
	"github.com/phper95/tinydocker/container"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "simple-runc",
		Usage: "A simple container runtime",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run a command in a new container",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "it",
						Usage: "Interactive mode with pseudo-TTY",
					},
				},
				Action: func(c *cli.Context) error {
					cmd := c.Args().Slice()
					if len(cmd) == 0 {
						log.Fatal("No command specified")
					}
					return container.Run(cmd, c.Bool("it"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
