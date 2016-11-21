package main

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/rmrobinson/hue-go"
)

var (
	bridgeAddr = flag.String("bridgeAddress", "", "The IP address of the bridge to connect to")
	username   = flag.String("username", "", "The username on the bridge to use")
)

func main() {
	flag.Parse()

	bridgeUrl, err := url.Parse("http://" + *bridgeAddr + "/description.xml")

	if err != nil {
		fmt.Printf("Unable to parse supplied bridge address: %s\n", err.Error())
		return
	}

	var b hue_go.Bridge

	err = b.Init(bridgeUrl)

	if err != nil {
		fmt.Printf("Unable to initialize supplied bridge address: %s\n", err.Error())
		return
	}

	b.Username = *username

	u := hue_go.NewUpdater(&b)

	results := make(chan string)
	quit := make(chan interface{})

	go u.Run(results, quit)

	for {
		msg := <-results

		fmt.Printf("Bridge updater changed: %s\n", msg)
	}

	return
}
