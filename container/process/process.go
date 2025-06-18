package process

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
)

func AttachTerminal(cmd *exec.Cmd) {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	defer ptmx.Close()

	go func() {
		ioCopy(ptmx, os.Stdout)
	}()
	go func() {
		ioCopy(os.Stdin, ptmx)
	}()
}

func ioCopy(dst *os.File, src *os.File) {
	defer dst.Close()
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil {
			return
		}
		dst.Write(buf[:n])
	}
}
