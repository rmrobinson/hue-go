package hue

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var (
	// ErrBridgeNotConfigured is returned if the bridge URL configured is not valid
	ErrBridgeNotConfigured = errors.New("bridge is not configured")
	// ErrInvalidManufacturer is returned if the retrieved bridge manufacturer is invalid
	ErrInvalidManufacturer = errors.New("invalid manufacturer detected")
	// ErrInvalidModel is returned if the retrieved bridge model is invalid
	ErrInvalidModel = errors.New("invalid model detected")
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
	ManufacturerURL  string `xml:"manufacturerURL"`
	ModelDescription string `xml:"modelDescription"`
	ModelName        string `xml:"modelName"`
	ModelNumber      string `xml:"modelNumber"`
	ModelURL         string `xml:"modelURL"`
	SerialNumber     string `xml:"serialNumber"`
	UDN              string `xml:"UDN"`
	PresentationURL  string `xml:"presentationURL"`
	Icons            []Icon `xml:"iconList>icon"`
}

// BridgeDescription contains the device information and the URL base which all subsequent calls should be made against.
type BridgeDescription struct {
	URLBase string `xml:"URLBase"`
	Device  Device `xml:"device"`
}

// Bridge represents an instance of a Hue bridge.
// The locator discovery process will also create one instance for each profile detected that is not already running.
type Bridge struct {
	id string

	// An empty username implies we have not yet paired with the bridge.
	Username string

	baseURL *url.URL

	validateURL *url.URL

	iconURL *url.URL

	updateInProgress bool
}

// NewBridge creates a new instance of a Hue bridge.
func NewBridge(username string) *Bridge {
	return &Bridge{
		Username: username,
	}
}

// InitURL initializes this bridge instance with the specified discovery URL.
func (b *Bridge) InitURL(validateURL *url.URL) error {
	b.validateURL = validateURL

	desc, err := b.Description()
	if err != nil {
		return err
	}

	b.id = desc.Device.SerialNumber
	b.baseURL, err = url.Parse(desc.URLBase)

	return err
}

// InitIP initializes this bridge instance with the specified IP address.
func (b *Bridge) InitIP(bridgeIP string) error {
	bridgeURL, err := bridgeDescURLFromIP(bridgeIP)
	if err != nil {
		return err
	}

	return b.InitURL(bridgeURL)
}

// Description parses the validation XML file present on every Hue bridge.
// See http://www.developers.meethue.com/documentation/hue-bridge-discovery for details of the response format.
func (b *Bridge) Description() (BridgeDescription, error) {
	desc := BridgeDescription{}

	if b.validateURL == nil {
		return desc, ErrBridgeNotConfigured
	}

	res, err := http.Get(b.validateURL.String())
	if err != nil {
		return desc, err
	}

	err = xml.NewDecoder(res.Body).Decode(&desc)
	if err != nil {
		return desc, err
	}

	if desc.Device.Manufacturer != "Royal Philips Electronics" {
		return desc, ErrInvalidManufacturer
	} else if !strings.Contains(desc.Device.ModelName, "hue bridge") {
		return desc, ErrInvalidModel
	}

	return desc, nil
}

func (b *Bridge) isAvailable() bool {
	return len(b.Username) > 0 && len(b.baseURL.Host) > 0
}

// ID returns the unique ID of the bridge
func (b *Bridge) ID() string {
	return b.id
}

// IsUpdating returns whether a bridge is in the process of updating or not.
func (b *Bridge) IsUpdating() bool {
	return b.updateInProgress
}
