package hue_go

type LightArg arg

func (l *LightArg) Reset() {
	l.args = make(map[string]interface{})
}

func (l *LightArg) Errors() map[string]responseError {
	return l.errors
}

func (l *LightArg) SetName(name string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["name"] = name
}

func (l *LightArg) Name() string {
	if ret, ok := l.args["name"].(string); ok {
		return ret
	} else {
		return ""
	}
}

type LightStateArg arg

func (l *LightStateArg) Reset() {
	l.args = make(map[string]interface{})
}

func (l *LightStateArg) Errors() map[string]responseError {
	return l.errors
}

func (l *LightStateArg) SetIsOn(isOn bool) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["on"] = isOn
}

func (l *LightStateArg) IsOn() bool {
	if ret, ok := l.args["on"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (l *LightStateArg) SetBrightness(brightness uint8) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["bri"] = brightness
}

func (l *LightStateArg) Brightness() uint8 {
	if ret, ok := l.args["bri"].(uint8); ok {
		return ret
	} else {
		return 0
	}
}

func (l *LightStateArg) SetHue(hue uint16) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["hue"] = hue
}

func (l *LightStateArg) Hue() uint16 {
	if ret, ok := l.args["hue"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

func (l *LightStateArg) SetSaturation(saturation uint8) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["sat"] = saturation
}

func (l *LightStateArg) Saturation() uint8 {
	if ret, ok := l.args["sat"].(uint8); ok {
		return ret
	} else {
		return 0
	}
}

func (l *LightStateArg) SetXY(xy XY) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["xy"] = [2]float64{xy.X, xy.Y}
}

func (l *LightStateArg) XY() XY {
	if ret, ok := l.args["xy"].([2]float64); ok {
		return XY{X: ret[0], Y: ret[1]}
	} else {
		return XY{}
	}
}

func (l *LightStateArg) SetColourTemperature(ct uint16) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["ct"] = ct
}

func (l *LightStateArg) ColourTemperature() uint16 {
	if ret, ok := l.args["ct"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

func (l *LightStateArg) SetRGB(rgb RGB, lightModel string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	var xy XY
	xy.FromRGB(rgb, lightModel)
	l.args["xy"] = [2]float64{xy.X, xy.Y}
}

func (l *LightStateArg) RGB(lightModel string) RGB {
	if ret, ok := l.args["xy"].([2]float64); ok {
		var rgb RGB
		rgb.FromXY(XY{X: ret[0], Y: ret[1]}, lightModel)
		return rgb
	} else {
		return RGB{}
	}
}

func (l *LightStateArg) SetTransitionTime(tt uint16) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["transitiontime"] = tt
}

func (l *LightStateArg) TransitionTime() uint16 {
	if ret, ok := l.args["transitiontime"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

func (l *LightStateArg) SetAlert(alert string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["alert"] = alert
}

func (l *LightStateArg) Alert() string {
	if ret, ok := l.args["alert"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (l *LightStateArg) SetEffect(effect string) {
	if l.args == nil {
		l.args = make(map[string]interface{})
	}

	l.args["effect"] = effect
}

func (l *LightStateArg) Effect() string {
	if ret, ok := l.args["effect"].(string); ok {
		return ret
	} else {
		return ""
	}
}
