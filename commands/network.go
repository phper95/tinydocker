package commands

import (
	"errors"
	"github.com/phper95/tinydocker/network"
	"github.com/urfave/cli"
)

// docker network create --subnet <cidr> --driver <driver> <name>
// docker network ls
// docker network rm <name>
var NetworkCommand = cli.Command{
	Name:  "network",
	Usage: "Manage networks",
	Subcommands: []cli.Command{
		{
			Name:  "create",
			Usage: "Create a network",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "subnet",
					Usage: "Subnet in CIDR format (e.g., 192.168.243.0/24)",
				},
				&cli.StringFlag{
					Name:  "driver",
					Usage: "Network driver (e.g., bridge)",
					Value: "bridge",
				},
			},
			Action: func(ctx *cli.Context) error {
				name := ctx.Args().First()
				if name == "" {
					return errors.New("network name cannot be empty")
				}
				subnet := ctx.String("subnet")
				if subnet == "" {
					return errors.New("--subnet is required")
				}
				driver := ctx.String("driver")
				return network.CreateNetwork(name, driver, subnet)
			},
		},
		{
			Name:  "ls",
			Usage: "List container network",
			Action: func(ctx *cli.Context) error {
				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "rm",
			Usage: "Remove one or more networks",
			Action: func(ctx *cli.Context) error {
				if ctx.NArg() < 1 {
					return errors.New("network name is required")
				}
				name := ctx.Args().First()
				if name == "" {
					return errors.New("network name cannot be empty")
				}
				return network.DeleteNetwork(name)
			},
		},
	},
}
