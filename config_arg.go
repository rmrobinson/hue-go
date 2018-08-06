package hue

// ConfigArg represents a config property that can be set.
type ConfigArg arg

// Reset clears this configuration.
func (c *ConfigArg) Reset() {
	c.args = make(map[string]interface{})
}

// Errors exposes any errors encountered when applying the configuration.
func (c *ConfigArg) Errors() map[string]ResponseError {
	return c.errors
}

// SetProxyPort saves the specified value to be applied.
func (c *ConfigArg) SetProxyPort(port uint16) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["proxyport"] = port
}

// ProxyPort returns the proxy port option, if configured.
func (c *ConfigArg) ProxyPort() uint16 {
	if ret, ok := c.args["proxyport"].(uint16); ok {
		return ret
	}
	return 0
}

// SetProxyAddress saves the specified value to be applied.
func (c *ConfigArg) SetProxyAddress(address string) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	if len(address) < 1 {
		c.args["proxyaddress"] = "none"
	} else {
		c.args["proxyaddress"] = address
	}
}

// ProxyAddress returns the proxy address option, if configured.
func (c *ConfigArg) ProxyAddress() string {
	if ret, ok := c.args["proxyaddress"].(string); ok {
		return ret
	}
	return ""
}

// SetName saves the specified value to be applied.
func (c *ConfigArg) SetName(name string) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["name"] = name
}

// Name returns the name option, if configured.
func (c *ConfigArg) Name() string {
	if ret, ok := c.args["name"].(string); ok {
		return ret
	}
	return ""
}

// SetDHCP saves the specified value to be applied.
func (c *ConfigArg) SetDHCP() {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["dhcp"] = true
}

// DHCP returns the DHCP option, if configured.
func (c *ConfigArg) DHCP() bool {
	if ret, ok := c.args["dhcp"].(bool); ok {
		return ret
	}
	return false
}

// SetStaticAddress saves the specified value to be applied.
func (c *ConfigArg) SetStaticAddress(addr string, netmask string, gateway string) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["ipaddress"] = addr
	c.args["netmask"] = netmask
	c.args["gateway"] = gateway
	c.args["dhcp"] = false
}

// StaticAddress returns the static address option, if configured.
func (c *ConfigArg) StaticAddress() (string, string, string) {
	var addr, netmask, gateway string
	var ok bool

	addr, ok = c.args["addr"].(string)

	if !ok {
		return "", "", ""
	}

	netmask, ok = c.args["netmask"].(string)

	if !ok {
		return "", "", ""
	}

	gateway, ok = c.args["gateway"].(string)

	if !ok {
		return "", "", ""
	}

	return addr, netmask, gateway
}

// SetTouchLink saves the specified value to be applied.
func (c *ConfigArg) SetTouchLink() {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["touchlink"] = true
}
