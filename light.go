package hue_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

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

type Light struct {
	Id               string
	Name             string `json:"name"`
	Model            string `json:"type"`
	ModelId          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	UniqueId         string `json:"uniqueid"`
	SwVersion        string `json:"swversion"`

	State LightState `json:"state"`
}

type NewLight struct {
	Id   string
	Name string `json:"name"`
}

func (b *Bridge) SearchForNewLights() error {
	return errors.New("Search is not yet supported")
}

func (b *Bridge) NewLights() (lights []NewLight, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/lights/new"

	resp, err := http.Get(url)

	var respEntries map[string]*json.RawMessage

	err = json.NewDecoder(resp.Body).Decode(&respEntries)

	for key, respEntry := range respEntries {
		if key == "lastscan" {
			var time string
			if err = json.Unmarshal(*respEntry, &time); err != nil {
				return
			}

			// TODO: do something with this value?
		} else {
			var l NewLight
			if err = json.Unmarshal(*respEntry, &l); err != nil {
				return
			}

			l.Id = key

			lights = append(lights, l)
		}
	}

	return
}

func (b *Bridge) Lights() (lights []Light, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/lights"

	res, err := http.Get(url)

	var respBody map[string]Light

	err = json.NewDecoder(res.Body).Decode(&respBody)

	if err != nil {
		return
	}

	for id, light := range respBody {
		light.Id = id

		if light.State.ColorMode == "xy" {
			xy := XY{X: light.State.XY[0], Y: light.State.XY[1]}
			light.State.RGB.FromXY(xy, light.ModelId)
		} else if light.State.ColorMode == "ct" {
			light.State.RGB.FromCT(light.State.ColorTemperature)
		} else if light.State.ColorMode == "hs" {
			hsb := HSB{Hue: light.State.Hue, Saturation: light.State.Saturation, Brightness: light.State.Brightness}
			light.State.RGB.FromHSB(hsb)
		}

		lights = append(lights, light)
	}

	return
}

func (b *Bridge) Light(id string) (light Light, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/lights/" + id

	resp, err := http.Get(url)

	err = json.NewDecoder(resp.Body).Decode(&light)

	if err != nil {
		return
	}

	if light.State.ColorMode == "xy" {
		xy := XY{X: light.State.XY[0], Y: light.State.XY[1]}
		light.State.RGB.FromXY(xy, light.ModelId)
	} else if light.State.ColorMode == "ct" {
		light.State.RGB.FromCT(light.State.ColorTemperature)
	} else if light.State.ColorMode == "hs" {
		hsb := HSB{Hue: light.State.Hue, Saturation: light.State.Saturation, Brightness: light.State.Brightness}
		light.State.RGB.FromHSB(hsb)
	}

	return
}

func (b *Bridge) SetLight(id string, args *LightArg) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/lights/" + id

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(args.args)

	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	var respEntries responseEntries

	err = json.NewDecoder(resp.Body).Decode(&respEntries)

	for _, respEntry := range respEntries {
		var e responseEntry
		if err = json.Unmarshal(respEntry, &e); err != nil {
			return
		}

		if e.Error.Type > 0 {
			if args.errors == nil {
				args.errors = make(map[string]responseError)
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
						return
					}

					args.args[key] = v
				}
			}
		}
	}

	return
}

func (b *Bridge) SetLightState(id string, args *LightStateArg) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/lights/" + id + "/state"

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(args.args)

	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, buf)

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	var respEntries responseEntries

	err = json.NewDecoder(resp.Body).Decode(&respEntries)

	for _, respEntry := range respEntries {
		var e responseEntry
		if err = json.Unmarshal(respEntry, &e); err != nil {
			return
		}

		if e.Error.Type > 0 {
			if args.errors == nil {
				args.errors = make(map[string]responseError)
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
						return
					}

					args.args[key] = xyArr
				} else if key == "bri" || key == "hue" || key == "sat" || key == "ct" || key == "transitiontime" {
					var v int
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "on" {
					var v bool
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "alert" || key == "effect" {
					var v string
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				}
			}
		}
	}

	return
}

func (b *Bridge) DeleteLight(id string) error {
	return errors.New("Delete is not yet supported")
}
