//go:build !linux || !cgo

package container

import (
    "fmt"

    "github.com/urfave/cli"
)

func ExecContainer(args cli.Args) error {
    return fmt.Errorf("exec is only supported on linux with cgo enabled")
}

