package hue_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type SensorState struct {
	LastUpdated string `json:"lastupdated"`

	ButtonEvent int32 `json:"buttonevent"`

	Presence    bool  `json:"presence"`
	Temperature int32 `json:"temperature"`
	Humidity    int32 `json:"humidity"`

	Daylight bool `json:"daylight"`

	LightLevel uint16 `json:"lightlevel"`
	Dark       bool   `json:"dark"`

	Flag bool `json:"flag"`

	Status int32 `json:"status"`
}

type SensorConfig struct {
	On        bool `json:"on"`
	Reachable bool `json:"reachable"`

	Battery  uint8    `json:"battery"`
	Alert    string   `json:"alert"`
	UserTest bool     `json:"usertest"`
	Url      *url.URL `json:"url"`

	Pending       []string `json:"pending"`
	LedIndication bool     `json:"ledindication"`

	Latitude   string `json:"lat"`
	Longitude  string `json:"long"`
	Configured bool   `json:"configured"`

	SunriseOffset int8 `json:"sunriseoffset"`
	SunsetOffset  int8 `json:"sunsetoffset"`

	ThresholdDark   uint16 `json:"tholddark"`
	ThresholdOffset uint16 `json:"tholdoffset"`
}

type Sensor struct {
	Id               string
	Name             string `json:"name"`
	Type             string `json:"type"`
	ModelId          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	ProductId        string `json:"productid"`
	UniqueId         string `json:"uniqueid"`
	SwVersion        string `json:"swversion"`

	Sensitivity    int32 `json:"sensitivity"`
	SensitivityMax int32 `json:"sensitivitymax"`

	Config SensorConfig `json:"config"`

	State SensorState `json:"state"`
}

type NewSensor struct {
	Id   string
	Name string `json:"name"`
}

func (b *Bridge) SearchForNewSensors() error {
	return errors.New("Search is not yet supported")
}

func (b *Bridge) NewSensors() (sensors []NewSensor, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors/new"

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
			var s NewSensor
			if err = json.Unmarshal(*respEntry, &s); err != nil {
				return
			}

			s.Id = key

			sensors = append(sensors, s)
		}
	}

	return
}

func (b *Bridge) Sensors() (sensors []Sensor, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors"

	res, err := http.Get(url)

	var respBody map[string]Sensor

	err = json.NewDecoder(res.Body).Decode(&respBody)

	if err != nil {
		return
	}

	for id, sensor := range respBody {
		sensor.Id = id

		sensors = append(sensors, sensor)
	}

	return
}

func (b *Bridge) Sensor(id string) (sensor Sensor, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors/" + id

	resp, err := http.Get(url)

	err = json.NewDecoder(resp.Body).Decode(&sensor)

	if err != nil {
		return
	}

	return
}

func (b *Bridge) SetSensor(id string, args *SensorArg) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors/" + id

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

func (b *Bridge) SetSensorConfig(id string, args *SensorConfigArg) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors/" + id + "/config"

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

				if key == "alert" || key == "url" || key == "lat" || key == "long" {
					var v string
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "on" || key == "reachable" {
					var v bool
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "battery" {
					var v uint8
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "sunriseoffset" || key == "sunsetoffset" {
					var v int8
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "tholddark" || key == "tholdoffset" {
					var v uint16
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

func (b *Bridge) SetSensorState(id string, args *SensorStateArg) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors/" + id + "/state"

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

				if key == "open" || key == "presence" || key == "flag" {
					var v bool
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "temperature" || key == "humidity" || key == "status" {
					var v int32
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					args.args[key] = v
				} else if key == "lightlevel" {
					var v uint16
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

func (b *Bridge) CreateSensor(sensor *Sensor) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/sensors"

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(sensor)

	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)

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
			return errors.New(e.Error.Description)
		} else {
			for path, jsonValue := range e.Success {
				keys := strings.Split(path, "/")

				key := keys[len(keys)-1]

				if key == "id" {
					var v string
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return
					}

					sensor.Id = v
				}
			}
		}
	}

	return
}

func (b *Bridge) DeleteSensor(id string) error {
	return errors.New("Delete is not yet supported")
}
