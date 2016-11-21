package hue_go

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"log"

	"github.com/huin/goupnp/ssdp"
)

const (
	SOURCE_STATIC = iota
	SOURCE_UPNP
	SOURCE_NUPNP
)

type result struct {
	url    *url.URL
	source int
}

type Locator struct {
	staticAddrs []string

	// Map the ID to the validation URL last seen with that ID.
	profiles map[string]*url.URL

	incoming chan result
}

func NewLocator() *Locator {
	d := &Locator{
		incoming:    make(chan result),
		profiles:    make(map[string]*url.URL),
		staticAddrs: make([]string, 8),
	}

	return d
}

func (d *Locator) AddStaticAddress(addr string) {
	d.staticAddrs = append(d.staticAddrs, addr)
}

func (d *Locator) Run(results chan Bridge) {
	go runStatic(d.staticAddrs, d.incoming)
	go runUPnP(d.incoming)
	go runNUPnP(d.incoming)

	for {
		res := <-d.incoming

		if res.url == nil {
			continue
		}

		br := Bridge{}
		err := br.Init(res.url)

		if err != nil {
			log.Fatalf("Unable to validate bridge URL %s: %s\n", res.url.String(), err)
			continue
		}

		currUrl, ok := d.profiles[br.Id()]

		// If the ID isn't present, we haven't seen this bridge before
		if !ok {
			log.Printf("New record found (url = %s) via %d, reporting\n", res.url.String(), res.source)

			d.profiles[br.Id()] = res.url

			results <- br
		} else if !(strings.Contains(currUrl.Host, res.url.Host) || strings.Contains(res.url.Host, currUrl.Host)) || currUrl.Path != res.url.Path || currUrl.Scheme != res.url.Scheme {
			// We don't do a stright value comparison because UPnP returns the port, while nUPnP does not.
			// So we use the host comparisons to drop the port, then check path and scheme.
			// This could likely be improved, possibly by adding :80 in the nUPnP detection scheme.
			log.Printf("Bridge %s changed, new validation URL is %s (old was %s)\n", br.Id(), res.url.String(), currUrl.String())

			// TODO: update bridge

		}
		// We don't need to log when we find the same bridge over and over.
	}
}

func runStatic(addrs []string, results chan result) {
	for _, addr := range addrs {
		// Skip empty addresses
		if len(addr) < 1 {
			continue
		}

		url, _ := url.Parse("http://" + addr + "/description.xml")

		r := result{
			url:    url,
			source: SOURCE_STATIC,
		}
		results <- r
	}
}

func runNUPnP(results chan result) {
	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			res, err := http.Get("https://www.meethue.com/api/nupnp")

			if err != nil {
				continue
			}

			type entry struct {
				Id                 string `json:"id"`
				InternallIpAddress string `json:"internalipaddress"`
				MacAddress         string `json:"macaddress"`
				Name               string `json:"name"`
			}

			var body []entry

			err = json.NewDecoder(res.Body).Decode(&body)

			if err != nil {
				continue
			}

			for _, entry := range body {
				// From http://www.developers.meethue.com/documentation/hue-bridge-discovery
				// We assume that the bridge will always have an XML description file present
				// when the N-UPnP approach is used.
				url, err := url.Parse("http://" + entry.InternallIpAddress + "/description.xml")

				if err != nil {
					continue
				}

				r := result{
					url:    url,
					source: SOURCE_NUPNP,
				}
				results <- r
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func runUPnP(results chan result) {
	c := make(chan ssdp.Update)
	srv, reg := ssdp.NewServerAndRegistry()
	reg.AddListener(c)
	go runSSDPReceiver(c, results)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe failed: ", err)
	}
}

func runSSDPReceiver(c <-chan ssdp.Update, results chan result) {
	for u := range c {
		if u.Entry == nil {
			continue
		} else if !strings.Contains(u.Entry.Server, "IpBridge") {
			continue
		}

		r := result{
			url:    &u.Entry.Location,
			source: SOURCE_UPNP,
		}
		results <- r
	}
}
