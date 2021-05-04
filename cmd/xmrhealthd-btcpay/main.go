package main

import (
	"log"
	"os/exec"
	"strings"

	"github.com/dys2p/xmrhealthd"
)

func main() {

	var ipAddress string

	if ip, err := exec.Command("/usr/bin/docker", "inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", "btcpayserver_monerod").Output(); err == nil {
		ipAddress = strings.TrimSpace(string(ip))
	} else {
		log.Println(err)
		return
	}

	xmrhealthd.Run(ipAddress)
}
