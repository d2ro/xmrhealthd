package main

import (
	"os"

	"github.com/dys2p/xmrhealthd"
)

func main() {

	var ipAddress string

	if len(os.Args) > 1 {
		ipAddress = os.Args[1]
	} else {
		ipAddress = "127.0.0.1"
	}

	xmrhealthd.Run(ipAddress)
}
