package hue_go

type ConfigArg arg

func (c *ConfigArg) Reset() {
	c.args = make(map[string]interface{})
}

func (c *ConfigArg) Errors() map[string]responseError {
	return c.errors
}

func (c *ConfigArg) SetProxyPort(port uint16) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["proxyport"] = port
}

func (c *ConfigArg) ProxyPort() uint16 {
	if ret, ok := c.args["proxyport"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

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

func (c *ConfigArg) ProxyAddress() string {
	if ret, ok := c.args["proxyaddress"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (c *ConfigArg) SetName(name string) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["name"] = name
}

func (c *ConfigArg) Name() string {
	if ret, ok := c.args["name"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (c *ConfigArg) SetDhcp() {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["dhcp"] = true
}

func (c *ConfigArg) Dhcp() bool {
	if ret, ok := c.args["dhcp"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (c *ConfigArg) SetStaticAddress(addr string, netmask string, gateway string) {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["ipaddress"] = addr
	c.args["netmask"] = netmask
	c.args["gateway"] = gateway
	c.args["dhcp"] = false
}

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

func (c *ConfigArg) SetTouchLink() {
	if c.args == nil {
		c.args = make(map[string]interface{})
	}

	c.args["touchlink"] = true
}
