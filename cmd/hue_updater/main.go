package main

import (
	"flag"
	"fmt"
	"github.com/rmrobinson/hue-go"
)

var (
	bridgeAddr = flag.String("bridgeAddress", "", "The IP address of the bridge to connect to")
	username   = flag.String("username", "", "The username on the bridge to use")
)

func main() {
	flag.Parse()

	b := hue.NewBridge(*username)
	err := b.InitIP(*bridgeAddr)
	if err != nil {
		fmt.Printf("Unable to initialize supplied bridge address: %s\n", err.Error())
		return
	}

	u := hue.NewUpdater(b)

	results := make(chan string)
	quit := make(chan interface{})

	go u.Run(results, quit)

	for {
		msg := <-results

		fmt.Printf("Bridge updater changed: %s\n", msg)
	}

	return
}
