package local_pkg

import (
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func LocalFunc() {
	log.Println("Hello from local package")
}

//Delve
//go install github.com/go-delve/delve/cmd/dlv@latest
//dlv debug /path/to/binary
