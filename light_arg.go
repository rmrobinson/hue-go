package hue

// LightArg represents a configuration argument that can be made to a light.
type LightArg arg

// Reset clears any set configuration options,
func (l *LightArg) Reset() {
	l.args = make(map[string]interface{})
}

// Errors returns any errors encountered when applying the specified configuration.
func (l *LightArg) Errors() map[string]ResponseError {
	return l.errors
}

// SetName saves the specified value to be applied.
func (l *LightArg) SetName(name string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["name"] = name
}

// Name returns the name option, if set.
func (l *LightArg) Name() string {
	if ret, ok := l.args["name"].(string); ok {
		return ret
	}
	return ""
}

// LightStateArg represents a state setting that can be applied to a light.
type LightStateArg arg

// Reset clears anything set.
func (l *LightStateArg) Reset() {
	l.args = make(map[string]interface{})
}

// Errors returns any errors encountered when applying the specified setting changes.
func (l *LightStateArg) Errors() map[string]ResponseError {
	return l.errors
}

// SetIsOn saves the specified value to be applied.
func (l *LightStateArg) SetIsOn(isOn bool) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["on"] = isOn
}

// IsOn returns the 'is on' value, if configured.
func (l *LightStateArg) IsOn() bool {
	if ret, ok := l.args["on"].(bool); ok {
		return ret
	}
	return false
}

// SetBrightness saves the specified value to be applied.
func (l *LightStateArg) SetBrightness(brightness uint8) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["bri"] = brightness
}

// Brightness returns the brightness value, if configured.
func (l *LightStateArg) Brightness() uint8 {
	if ret, ok := l.args["bri"].(uint8); ok {
		return ret
	}
	return 0
}

// SetHue saves the specified value to be applied.
func (l *LightStateArg) SetHue(hue uint16) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["hue"] = hue
}

// Hue returns the hue value, if configured.
func (l *LightStateArg) Hue() uint16 {
	if ret, ok := l.args["hue"].(uint16); ok {
		return ret
	}
	return 0
}

// SetSaturation saves the specified value to be applied.
func (l *LightStateArg) SetSaturation(saturation uint8) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["sat"] = saturation
}

// Saturation returns the saturation value, if configured.
func (l *LightStateArg) Saturation() uint8 {
	if ret, ok := l.args["sat"].(uint8); ok {
		return ret
	}
	return 0
}

// SetXY saves the specified value to be applied.
func (l *LightStateArg) SetXY(xy XY) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["xy"] = [2]float64{xy.X, xy.Y}
}

// XY returns the XY value, if configured.
func (l *LightStateArg) XY() XY {
	if ret, ok := l.args["xy"].([2]float64); ok {
		return XY{X: ret[0], Y: ret[1]}
	}
	return XY{}
}

// SetColourTemperature saves the specified value to be applied.
func (l *LightStateArg) SetColourTemperature(ct uint16) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["ct"] = ct
}

// ColourTemperature returns the colour temperature value, if configured.
func (l *LightStateArg) ColourTemperature() uint16 {
	if ret, ok := l.args["ct"].(uint16); ok {
		return ret
	}
	return 0
}

// SetRGB saves the specified value to be applied.
func (l *LightStateArg) SetRGB(rgb RGB, lightModel string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	var xy XY
	xy.FromRGB(rgb, lightModel)
	l.args["xy"] = [2]float64{xy.X, xy.Y}
}

// RGB returns the RGB value, if configured.
func (l *LightStateArg) RGB(lightModel string) RGB {
	if ret, ok := l.args["xy"].([2]float64); ok {
		var rgb RGB
		rgb.FromXY(XY{X: ret[0], Y: ret[1]}, lightModel)
		return rgb
	}
	return RGB{}
}

// SetTransitionTime saves the specified value to be applied.
func (l *LightStateArg) SetTransitionTime(tt uint16) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["transitiontime"] = tt
}

// TransitionTime returns the transition time, if configured.
func (l *LightStateArg) TransitionTime() uint16 {
	if ret, ok := l.args["transitiontime"].(uint16); ok {
		return ret
	}
	return 0
}

// SetAlert saves the specified value to be applied.
func (l *LightStateArg) SetAlert(alert string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["alert"] = alert
}

// Alert returns the alert value, if configured.
func (l *LightStateArg) Alert() string {
	if ret, ok := l.args["alert"].(string); ok {
		return ret
	}
	return ""
}

// SetEffect saves the specified value to be applied.
func (l *LightStateArg) SetEffect(effect string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["effect"] = effect
}

// Effect returns the effect, if configured.
func (l *LightStateArg) Effect() string {
	if ret, ok := l.args["effect"].(string); ok {
		return ret
	}
	return ""
}
