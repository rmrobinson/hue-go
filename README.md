# hue-go

A golang library that wraps the Philips Hue REST API (API 1.15)

The library currently supports the following Hue API endpoints:
 * Lights
 * Sensors
 * Config

The library supports auto-detection of bridges using the Locator functionality.

The library supports auto-update of bridges using the Updater functionality.

## Usage

Create an instance of the hue.Bridge struct, then call InitIP() with the IP of the bridge.
If there is no error, set the Username property on the instance of the hue.Bridge to an account which is authorized to access the API.
To acquire a username to use, call Pair() after calling InitIP(). This will set the Username on the bridge instance with a new username. If you wish to save the username for future use, it can be accessed via the Username property on the bridge instance.. It is possible but not advisable to call Pair() on a bridge with a username already set.

Once a Username is set, the library maps function calls on a 1:1 basis with the Hue REST API. For example, Lights() calls GET /lights, SetLightState calls PUT/lights/<ID>/state, etc.

The Set functions take a corresponding Arg which specifies which of the properties are to be saved by POST or PUT calls. Each PUT endpoint has a corresponding Arg type, with the collection of valid properties exposed with Getters and Setters.
The supplied Arg object has the updated value of each property after calling the Set method. If any of the properties failed to set, the Errors() function returns the details of why that property was unable to be saved.

An example of this can be found in examples/hue

It is possible to automatically get initialized Bridge instances by running a copy of the hue.Locator object in it's own goroutine.
The locator has 3 possible ways to locate bridges:
 * Static addresses configured prior to Run() being called by calling AddStaticAddress (max 8 currently supported).
 * UPnP on the current network
 * N-UPnP

An example of this can be found in examples/hue_locator

It is possible to automatically update each Bridge by creating an instance of the hue.Updater object, then calling Run() in a goroutine.
This updater will poll the bridge once an hour to determine if there is an update available; if there is it will automatically apply the update then continue monitoring for future updates.

An example of this can be found in examples/hue_updater

## TODO

* Groups
* Schedules
* Scenes
* Rules
* Resource links
