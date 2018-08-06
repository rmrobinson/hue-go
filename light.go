package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// LightState represents the state of a single light.
type LightState struct {
	On               bool       `json:"on"`
	Brightness       uint8      `json:"bri"`
	Hue              uint16     `json:"hue"`
	Saturation       uint8      `json:"sat"`
	XY               [2]float64 `json:"xy"`
	ColorTemperature uint16     `json:"ct"`
	RGB              RGB
	Alert            string `json:"alert"`
	Effect           string `json:"effect"`
	ColorMode        string `json:"colormode"`
	Reachable        bool   `json:"reachable"`
}

// Light represents a single light device.
type Light struct {
	ID               string
	Name             string `json:"name"`
	Model            string `json:"type"`
	ModelID          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	UniqueID         string `json:"uniqueid"`
	SwVersion        string `json:"swversion"`

	State LightState `json:"state"`
}

// NewLight represents a single, un-configured light.
type NewLight struct {
	ID   string
	Name string `json:"name"`
}

// SearchForNewLights begins the process of locating new lights on the bridge.
func (b *Bridge) SearchForNewLights() error {
	return errors.New("search is not yet supported")
}

// NewLights returns the collection of newly discovered lights.
// Only returns a value if SearchForNewLights has previously been called.
func (b *Bridge) NewLights() ([]NewLight, error) {
	if !b.isAvailable() {
		return nil, ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return nil, ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/lights/new"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var respEntries map[string]*json.RawMessage

	err = json.NewDecoder(resp.Body).Decode(&respEntries)
	if err != nil {
		return nil, err
	}

	var lights []NewLight
	for key, respEntry := range respEntries {
		if key == "lastscan" {
			var time string
			if err = json.Unmarshal(*respEntry, &time); err != nil {
				return nil, err
			}

			// TODO: do something with this value?
		} else {
			var l NewLight
			if err = json.Unmarshal(*respEntry, &l); err != nil {
				return nil, err
			}

			l.ID = key

			lights = append(lights, l)
		}
	}

	return lights, nil
}

// Lights returns the collection of lights configured on the bridge.
func (b *Bridge) Lights() ([]Light, error) {
	if !b.isAvailable() {
		return nil, ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return nil, ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/lights"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var respBody map[string]Light

	err = json.NewDecoder(res.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	var lights []Light
	for id, light := range respBody {
		light.ID = id

		if light.State.ColorMode == "xy" {
			xy := XY{X: light.State.XY[0], Y: light.State.XY[1]}
			light.State.RGB.FromXY(xy, light.ModelID)
		} else if light.State.ColorMode == "ct" {
			light.State.RGB.FromCT(light.State.ColorTemperature)
		} else if light.State.ColorMode == "hs" {
			hsb := HSB{Hue: light.State.Hue, Saturation: light.State.Saturation, Brightness: light.State.Brightness}
			light.State.RGB.FromHSB(hsb)
		}

		lights = append(lights, light)
	}

	return lights, err
}

// Light returns a single light from the bridge.
func (b *Bridge) Light(id string) (Light, error) {
	if !b.isAvailable() {
		return Light{}, ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return Light{}, ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/lights/" + id

	resp, err := http.Get(url)
	if err != nil {
		return Light{}, err
	}

	var light Light
	err = json.NewDecoder(resp.Body).Decode(&light)
	if err != nil {
		return light, err
	}

	if light.State.ColorMode == "xy" {
		xy := XY{X: light.State.XY[0], Y: light.State.XY[1]}
		light.State.RGB.FromXY(xy, light.ModelID)
	} else if light.State.ColorMode == "ct" {
		light.State.RGB.FromCT(light.State.ColorTemperature)
	} else if light.State.ColorMode == "hs" {
		hsb := HSB{Hue: light.State.Hue, Saturation: light.State.Saturation, Brightness: light.State.Brightness}
		light.State.RGB.FromHSB(hsb)
	}

	return light, nil
}

// SetLight updates the specified light with the configuration supplied.
func (b *Bridge) SetLight(id string, args *LightArg) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/lights/" + id

	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(args.args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var respEntries responseEntries
	err = json.NewDecoder(resp.Body).Decode(&respEntries)
	if err != nil {
		return err
	}

	for _, respEntry := range respEntries {
		var e responseEntry
		if err = json.Unmarshal(respEntry, &e); err != nil {
			return err
		}

		if e.Error.Type > 0 {
			if args.errors == nil {
				args.errors = make(map[string]ResponseError)
			}

			keys := strings.Split(e.Error.Address, "/")
			key := keys[len(keys)-1]

			args.errors[key] = e.Error
		} else {
			for path, jsonValue := range e.Success {
				keys := strings.Split(path, "/")

				key := keys[len(keys)-1]

				if key == "name" {
					var v string
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				}
			}
		}
	}

	return nil
}

// SetLightState sets the specified light with the supplied light state.
func (b *Bridge) SetLightState(id string, args *LightStateArg) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/lights/" + id + "/state"

	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(args.args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var respEntries responseEntries
	err = json.NewDecoder(resp.Body).Decode(&respEntries)

	for _, respEntry := range respEntries {
		var e responseEntry
		if err = json.Unmarshal(respEntry, &e); err != nil {
			return err
		}

		if e.Error.Type > 0 {
			if args.errors == nil {
				args.errors = make(map[string]ResponseError)
			}

			keys := strings.Split(e.Error.Address, "/")
			key := keys[len(keys)-1]

			args.errors[key] = e.Error
		} else {
			for path, jsonValue := range e.Success {
				keys := strings.Split(path, "/")

				key := keys[len(keys)-1]

				if key == "xy" {
					var xyArr []float64
					if err = json.Unmarshal(*jsonValue, &xyArr); err != nil {
						return err
					}

					args.args[key] = xyArr
				} else if key == "bri" || key == "hue" || key == "sat" || key == "ct" || key == "transitiontime" {
					var v int
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "on" {
					var v bool
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "alert" || key == "effect" {
					var v string
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				}
			}
		}
	}

	return nil
}

// DeleteLight removes the specified light from the bridge.
func (b *Bridge) DeleteLight(id string) error {
	return errors.New("delete is not yet supported")
}
