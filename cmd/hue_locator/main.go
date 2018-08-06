package main

import (
	"fmt"
	"github.com/rmrobinson/hue-go"
)

func main() {
	l := hue.NewLocator()
	results := make(chan hue.Bridge)

	go l.Run(results)

	for {
		b := <-results

		fmt.Printf("Bridge %s detected\n", b.ID())

		desc, err := b.Description()

		if err != nil {
			fmt.Printf("Unable to get description for bridge %s: %s\n", b.ID(), err.Error())
		} else {
			fmt.Printf("Bridge desc: %+v\n", desc)
		}
	}

	return
}
