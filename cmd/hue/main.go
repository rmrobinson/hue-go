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

	c, err := b.Config()

	if err != nil {
		fmt.Printf("Unable to retrieve bridge config: %s\n", err.Error())
		return
	}

	fmt.Printf("Bridge config: %+v\n", c)

	lights, err := b.Lights()

	if err != nil {
		fmt.Printf("Unable to retrieve bridge lights: %s\n", err.Error())
		return
	}

	for _, light := range lights {
		fmt.Printf("Light: %+v\n", light)
	}

	newLights, err := b.NewLights()

	if err != nil {
		fmt.Printf("Unable to retrieve new bridge lights: %s\n", err.Error())
		return
	}

	for _, newLight := range newLights {
		fmt.Printf("New Light: %+v\n", newLight)
	}

	sensors, err := b.Sensors()

	if err != nil {
		fmt.Printf("Unable to retrieve bridge sensors: %s\n", err.Error())
		return
	}

	for _, sensor := range sensors {
		fmt.Printf("Sensor: %+v\n", sensor)
	}

	newSensors, err := b.NewSensors()

	if err != nil {
		fmt.Printf("Unable to retrieve new bridge sensors: %s\n", err.Error())
		return
	}

	for _, newSensor := range newSensors {
		fmt.Printf("New Sensor: %+v\n", newSensor)
	}

	/*
		var lsa hue_go.LightStateArg
		lsa.SetIsOn(false)

		b.SetLightState("1", &lsa)

		if len(lsa.Errors()) > 0 {
			fmt.Printf("Unable to set light state to on\n")

			for _, err := range lsa.Errors() {
				fmt.Printf("Error: %s\n", err.Description)
			}
		} else {
			fmt.Printf("Is light 1 on? %t\n", lsa.IsOn())
		}
	*/

	return
}
