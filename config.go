package hue_go

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Config struct {
	Name string `json:"name"`
	Id   string `json:"bridgeid"`

	SwVersion    string `json:"swversion"`
	ApiVersion   string `json:"apiversion"`
	ModelVersion string `json:"modelid"`

	LinkButton     bool   `json:"linkbutton"`
	IpAddress      string `json:"ipaddress"`
	MacAddress     string `json:"mac"`
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

func (b *Bridge) Config() (config Config, err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/config"

	res, err := http.Get(url)

	err = json.NewDecoder(res.Body).Decode(&config)

	return
}

func (b *Bridge) SetConfig(args *ConfigArg) (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	} else if b.updateInProgress {
		err = errors.New("Bridge is being updated")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/config"

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

func (b *Bridge) Pair(identifier string) (username string, err error) {
	url := b.baseUrl.String() + "api"

	type reqBody struct {
		DeviceType string `json:"devicetype"`
	}

	jsonReq := &reqBody{DeviceType: "hue-go#" + identifier}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(jsonReq)

	res, err := http.Post(url, "application/json; charset=utf-8", buf)

	if err != nil {
		return
	}

	var respBody struct {
		successData struct {
			username string
		} `json:"success,omitempty"`
	}

	err = json.NewDecoder(res.Body).Decode(&respBody)

	if err != nil {
		return
	}

	username = respBody.successData.username
	b.Username = respBody.successData.username

	return
}

func (b *Bridge) CheckForUpdate() (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/config"

	var reqBody struct {
		SwUpdate struct {
			CheckForUpdate bool `json:"checkforupdate"`
		} `json:"swupdate"`
	}

	reqBody.SwUpdate.CheckForUpdate = true

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(reqBody)

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
			err = errors.New(e.Error.Description)
			return
		}
	}

	return
}

func (b *Bridge) StartUpdate() (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/config"

	var reqBody struct {
		SwUpdate struct {
			UpdateState int `json:"updatestate"`
		} `json:"swupdate"`
	}

	reqBody.SwUpdate.UpdateState = 3

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(reqBody)

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
			return
		}

		if e.Error.Type > 0 {
			return errors.New(e.Error.Description)
		}
	}

	return nil
}

func (b *Bridge) FinishUpdate() (err error) {
	if !b.isAvailable() {
		err = errors.New("Bridge is not yet ready")
		return
	}

	url := b.baseUrl.String() + "api/" + b.Username + "/config"

	var reqBody struct {
		SwUpdate struct {
			Notify bool `json:"notify"`
		} `json:"swupdate"`
	}

	reqBody.SwUpdate.Notify = false

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(reqBody)

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
			return
		}

		if e.Error.Type > 0 {
			return errors.New(e.Error.Description)
		}
	}

	return nil
}
