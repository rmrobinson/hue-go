package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// SensorState represents the state of a sensor.
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

// SensorConfig represents the configuration of a sensor.
type SensorConfig struct {
	On        bool `json:"on"`
	Reachable bool `json:"reachable"`

	Battery  uint8    `json:"battery"`
	Alert    string   `json:"alert"`
	UserTest bool     `json:"usertest"`
	URL      *url.URL `json:"url"`

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

// Sensor represents a single sensor instance on a Hue bridge.
type Sensor struct {
	ID               string
	Name             string `json:"name"`
	Type             string `json:"type"`
	ModelID          string `json:"modelid"`
	ManufacturerName string `json:"manufacturername"`
	ProductID        string `json:"productid"`
	UniqueID         string `json:"uniqueid"`
	SwVersion        string `json:"swversion"`

	Sensitivity    int32 `json:"sensitivity"`
	SensitivityMax int32 `json:"sensitivitymax"`

	Config SensorConfig `json:"config"`

	State SensorState `json:"state"`
}

// NewSensor represents an instance of a newly configured sensor instance paired to a Hue bridge.
type NewSensor struct {
	ID   string
	Name string `json:"name"`
}

// SearchForNewSensors attempts to discover new sensors added to the bridge.
func (b *Bridge) SearchForNewSensors() error {
	return errors.New("search is not yet supported")
}

// NewSensors returns the list of newly configured sensors.
// This will not return anything if SearchForNewSensors has not previously been called.
func (b *Bridge) NewSensors() ([]NewSensor, error) {
	if !b.isAvailable() {
		return nil, ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return nil, ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors/new"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var respEntries map[string]*json.RawMessage
	err = json.NewDecoder(resp.Body).Decode(&respEntries)
	if err != nil {
		return nil, err
	}

	var sensors []NewSensor
	for key, respEntry := range respEntries {
		if key == "lastscan" {
			var time string
			if err = json.Unmarshal(*respEntry, &time); err != nil {
				return nil, err
			}

			// TODO: do something with this value?
		} else {
			var s NewSensor
			if err = json.Unmarshal(*respEntry, &s); err != nil {
				return nil, err
			}

			s.ID = key

			sensors = append(sensors, s)
		}
	}

	return sensors, nil
}

// Sensors returns the list of sensors available on the bridge.
func (b *Bridge) Sensors() ([]Sensor, error) {
	if !b.isAvailable() {
		return nil, ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return nil, ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var respBody map[string]Sensor
	err = json.NewDecoder(res.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	var sensors []Sensor
	for id, sensor := range respBody {
		sensor.ID = id

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

// Sensor returns the specified sensor.
func (b *Bridge) Sensor(id string) (Sensor, error) {
	var sensor Sensor

	if !b.isAvailable() {
		return sensor, ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return sensor, ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors/" + id

	resp, err := http.Get(url)
	if err != nil {
		return sensor, err
	}

	err = json.NewDecoder(resp.Body).Decode(&sensor)
	return sensor, err
}

// SetSensor updates the configuration of the specified sensor.
func (b *Bridge) SetSensor(id string, args *SensorArg) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors/" + id

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

// SetSensorConfig updates the configuration a sensor.
func (b *Bridge) SetSensorConfig(id string, args *SensorConfigArg) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors/" + id + "/config"

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

				if key == "alert" || key == "url" || key == "lat" || key == "long" {
					var v string
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "on" || key == "reachable" {
					var v bool
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "battery" {
					var v uint8
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "sunriseoffset" || key == "sunsetoffset" {
					var v int8
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "tholddark" || key == "tholdoffset" {
					var v uint16
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

// SetSensorState updates the state of the specified sensor.
func (b *Bridge) SetSensorState(id string, args *SensorStateArg) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors/" + id + "/state"

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

				if key == "open" || key == "presence" || key == "flag" {
					var v bool
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "temperature" || key == "humidity" || key == "status" {
					var v int32
					if err = json.Unmarshal(*jsonValue, &v); err != nil {
						return err
					}

					args.args[key] = v
				} else if key == "lightlevel" {
					var v uint16
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

// CreateSensor adds a new sensor to the bridge.
func (b *Bridge) CreateSensor(sensor *Sensor) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/sensors"

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(sensor)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
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
			return errors.New(e.Error.Description)
		}

		for path, jsonValue := range e.Success {
			keys := strings.Split(path, "/")

			key := keys[len(keys)-1]

			if key == "id" {
				var v string
				if err = json.Unmarshal(*jsonValue, &v); err != nil {
					return err
				}

				sensor.ID = v
			}
		}
	}

	return nil
}

// DeleteSensor removes the specified sensor from the bridge.
func (b *Bridge) DeleteSensor(id string) error {
	return errors.New("delete is not yet supported")
}
