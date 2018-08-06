package hue

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var (
	// ErrBridgeNotAvailable is returned if the specified bridge is not yet configured for access.
	ErrBridgeNotAvailable = errors.New("bridge is not yet ready")
	// ErrBridgeUpdating is returned if the specified bridge is currently being updated.
	ErrBridgeUpdating = errors.New("bridge is currently being updated")
)

// Config represents an instance of a Hue bridge's configuration.
type Config struct {
	Name string `json:"name"`
	ID   string `json:"bridgeid"`

	SwVersion    string `json:"swversion"`
	APIVersion   string `json:"apiversion"`
	ModelVersion string `json:"modelid"`

	LinkButton     bool   `json:"linkbutton"`
	IPAddress      string `json:"ipaddress"`
	MACAddress     string `json:"mac"`
	SubnetMask     string `json:"netmask"`
	GatewayAddress string `json:"gateway"`

	IsDhcpAcquired bool `json:"dhcp"`
	PortalServices bool `json:"portalservices"`

	ZigbeeChannel int32 `json:"zigbeechannel"`

	Timezone string `json:"timezone"`

	SwUpdate struct {
		CheckForUpdates bool   `json:"checkforupdates"`
		UpdateDetails   string `json:"url"`
		UpdateSummary   string `json:"text"`
		NotifyUser      bool   `json:"notify"`
		State           int32  `json:"updateState"`

		DeviceTypes struct {
			Bridge  bool     `json:"bridge"`
			Lights  []string `json:"lights"`
			Sensors []string `json:"sensors"`
		} `json:"devicetypes"`
	} `json:"swupdate"`

	PortalState struct {
		SignedOn      bool   `json:"signedon"`
		Incoming      bool   `json:"incoming"`
		Outgoing      bool   `json:"outgoing"`
		Communication string `json:"communication"`
	} `json:"portalstate"`
}

// Config returns the configuration of the bridge
func (b *Bridge) Config() (Config, error) {
	config := Config{}
	if !b.isAvailable() {
		return config, ErrBridgeNotAvailable
	}

	url := b.baseURL.String() + "api/" + b.Username + "/config"

	res, err := http.Get(url)
	if err != nil {
		return config, err
	}

	err = json.NewDecoder(res.Body).Decode(&config)
	return config, err
}

// SetConfig applies the specified config options to the bridge.
func (b *Bridge) SetConfig(args *ConfigArg) error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	} else if b.updateInProgress {
		return ErrBridgeUpdating
	}

	url := b.baseURL.String() + "api/" + b.Username + "/config"
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

// Pair sets up the bridge with a new user.
func (b *Bridge) Pair(appName string, identifier string) error {
	url := b.baseURL.String() + "api"

	type reqBody struct {
		DeviceType string `json:"devicetype"`
	}

	jsonReq := &reqBody{DeviceType: appName + "#" + identifier}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(jsonReq)

	res, err := http.Post(url, "application/json; charset=utf-8", buf)
	if err != nil {
		return err
	}

	var respBody struct {
		successData struct {
			username string
		} `json:"success,omitempty"`
	}

	err = json.NewDecoder(res.Body).Decode(&respBody)
	if err != nil {
		return err
	}

	b.Username = respBody.successData.username

	return nil
}

// CheckForUpdate returns whether there is a software update available for the Hue bridge.
func (b *Bridge) CheckForUpdate() error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	}

	url := b.baseURL.String() + "api/" + b.Username + "/config"

	var reqBody struct {
		SwUpdate struct {
			CheckForUpdate bool `json:"checkforupdate"`
		} `json:"swupdate"`
	}

	reqBody.SwUpdate.CheckForUpdate = true

	buf := new(bytes.Buffer)

	err := json.NewEncoder(buf).Encode(reqBody)
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
			err = errors.New(e.Error.Description)
			return err
		}
	}

	return nil
}

// StartUpdate kicks off the update process for the Hue bridge.
func (b *Bridge) StartUpdate() error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	}

	url := b.baseURL.String() + "api/" + b.Username + "/config"

	var reqBody struct {
		SwUpdate struct {
			UpdateState int `json:"updatestate"`
		} `json:"swupdate"`
	}

	reqBody.SwUpdate.UpdateState = 3

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reqBody)
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
			return errors.New(e.Error.Description)
		}
	}

	return nil
}

// FinishUpdate completes the update process
func (b *Bridge) FinishUpdate() error {
	if !b.isAvailable() {
		return ErrBridgeNotAvailable
	}

	url := b.baseURL.String() + "api/" + b.Username + "/config"

	var reqBody struct {
		SwUpdate struct {
			Notify bool `json:"notify"`
		} `json:"swupdate"`
	}

	reqBody.SwUpdate.Notify = false

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reqBody)
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
			return errors.New(e.Error.Description)
		}
	}

	return nil
}
