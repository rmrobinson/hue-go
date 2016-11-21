package hue_go

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// Icon contains the details of an icon hosted on the bridge.
type Icon struct {
	MimeType string `xml:"mimetype"`
	Height   int    `xml:"height"`
	Width    int    `xml:"width"`
	Depth    int    `xml:"depth"`
	FileName string `xml:"url"`
}

// Device contains the details of the bridge itself. These are all read-only properties.
type Device struct {
	DeviceType       string `xml:"deviceType"`
	FriendlyName     string `xml:"friendlyName"`
	Manufacturer     string `xml:"manufacturer"`
	ManufacturerUrl  string `xml:"manufacturerURL"`
	ModelDescription string `xml:"modelDescription"`
	ModelName        string `xml:"modelName"`
	ModelNumber      string `xml:"modelNumber"`
	ModelUrl         string `xml:"modelURL"`
	SerialNumber     string `xml:"serialNumber"`
	UDN              string `xml:"UDN"`
	PresentationUrl  string `xml:"presentationURL"`
	Icons            []Icon `xml:"iconList>icon"`
}

// BridgeDescription contains the device information and the URL base which all subsequent calls should be made against.
type BridgeDescription struct {
	UrlBase string `xml:"URLBase"`
	Device  Device `xml:"device"`
}

// Bridge represents an instance of a Hue bridge.
// The locator discovery process will also create one instance for each profile detected that is not already running.
type Bridge struct {
	id string

	// An empty username implies we have not yet paired with the bridge.
	Username string

	baseUrl *url.URL

	validateUrl *url.URL

	iconUrl *url.URL

	updateInProgress bool
}

func (b *Bridge) Init(validateUrl *url.URL) (err error) {
	b.validateUrl = validateUrl

	desc, err := b.Description()

	if err != nil {
		return
	}

	b.id = desc.Device.SerialNumber
	b.baseUrl, err = url.Parse(desc.UrlBase)

	return
}

// Parse the validation XML file present on every Hue bridge.
// See http://www.developers.meethue.com/documentation/hue-bridge-discovery for details of the response format.
func (b *Bridge) Description() (desc BridgeDescription, err error) {
	if b.validateUrl == nil {
		err = errors.New("Bridge is not configured")
		return
	}

	res, err := http.Get(b.validateUrl.String())

	if err != nil {
		return
	}

	err = xml.NewDecoder(res.Body).Decode(&desc)

	if desc.Device.Manufacturer != "Royal Philips Electronics" {
		err = errors.New("Invalid manufacturer detected")
		return
	} else if !strings.Contains(desc.Device.ModelName, "hue bridge") {
		err = errors.New("Invalid model name detected")
		return
	}

	return
}

func (b *Bridge) isAvailable() bool {
	return len(b.Username) > 0 && len(b.baseUrl.Host) > 0
}

func (b *Bridge) Id() string {
	return b.id
}

func (b *Bridge) IsUpdating() bool {
	return b.updateInProgress
}
