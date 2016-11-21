package hue_go

type SensorArg arg

func (s *SensorArg) Reset() {
	s.args = make(map[string]interface{})
}

func (s *SensorArg) Errors() map[string]responseError {
	return s.errors
}

func (s *SensorArg) SetName(name string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["name"] = name
}

func (s *SensorArg) Name() string {
	if ret, ok := s.args["name"].(string); ok {
		return ret
	} else {
		return ""
	}
}

type SensorConfigArg arg

func (s *SensorConfigArg) Reset() {
	s.args = make(map[string]interface{})
}

func (s *SensorConfigArg) Errors() map[string]responseError {
	return s.errors
}

func (s *SensorConfigArg) SetIsOn(isOn bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["on"] = isOn
}

func (s *SensorConfigArg) IsOn() bool {
	if ret, ok := s.args["on"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (s *SensorConfigArg) SetIsReachable(isReachable bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["reachable"] = isReachable
}

func (s *SensorConfigArg) IsReachable() bool {
	if ret, ok := s.args["isReachable"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (s *SensorConfigArg) SetBatteryLevel(level uint8) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["battery"] = level
}

func (s *SensorConfigArg) BatteryLevel() uint8 {
	if ret, ok := s.args["battery"].(uint8); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorConfigArg) SetAlert(alert string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["alert"] = alert
}

func (s *SensorConfigArg) Alert() string {
	if ret, ok := s.args["alert"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (s *SensorConfigArg) SetURL(url string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["url"] = url
}

func (s *SensorConfigArg) Url() string {
	if ret, ok := s.args["url"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (s *SensorConfigArg) SetLatitude(latitude string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["lat"] = latitude
}

func (s *SensorConfigArg) Latitude() string {
	if ret, ok := s.args["lat"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (s *SensorConfigArg) SetLongitude(longitude string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["long"] = longitude
}

func (s *SensorConfigArg) Longitude() string {
	if ret, ok := s.args["long"].(string); ok {
		return ret
	} else {
		return ""
	}
}

func (s *SensorConfigArg) SetSunriseOffset(offset int8) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["sunriseoffset"] = offset
}

func (s *SensorConfigArg) SunriseOffset() int8 {
	if ret, ok := s.args["sunriseoffset"].(int8); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorConfigArg) SetSunsetOffset(offset int8) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["sunsetoffset"] = offset
}

func (s *SensorConfigArg) SunsetOffset() int8 {
	if ret, ok := s.args["sunsetoffset"].(int8); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorConfigArg) SetDarkThreshold(threshold uint16) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["tholddark"] = threshold
}

func (s *SensorConfigArg) DarkThreshold() uint16 {
	if ret, ok := s.args["tholddark"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorConfigArg) SetDarkThresholdOffset(offset uint16) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["tholdoffset"] = offset
}

func (s *SensorConfigArg) DarkThresholdOffset() uint16 {
	if ret, ok := s.args["tholdoffset"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

type SensorStateArg arg

func (s *SensorStateArg) Reset() {
	s.args = make(map[string]interface{})
}

func (s *SensorStateArg) Errors() map[string]responseError {
	return s.errors
}

func (s *SensorStateArg) SetIsOpen(isOpen bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["open"] = isOpen
}

func (s *SensorStateArg) IsOpen() bool {
	if ret, ok := s.args["open"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (s *SensorStateArg) SetIsPresent(isPresent bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["presence"] = isPresent
}

func (s *SensorStateArg) IsPresent() bool {
	if ret, ok := s.args["presence"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (s *SensorStateArg) SetTemperature(temperature int32) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["temperature"] = temperature
}

func (s *SensorStateArg) Temperature() int32 {
	if ret, ok := s.args["temperature"].(int32); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorStateArg) SetHumidity(humidity int32) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["humidity"] = humidity
}

func (s *SensorStateArg) Humidity() int32 {
	if ret, ok := s.args["humidity"].(int32); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorStateArg) SetLightLevel(lightLevel uint16) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["lightlevel"] = lightLevel
}

func (s *SensorStateArg) LightLevel() uint16 {
	if ret, ok := s.args["lightlevel"].(uint16); ok {
		return ret
	} else {
		return 0
	}
}

func (s *SensorStateArg) SetFlag(flag bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["flag"] = flag
}

func (s *SensorStateArg) Flag() bool {
	if ret, ok := s.args["flag"].(bool); ok {
		return ret
	} else {
		return false
	}
}

func (s *SensorStateArg) SetStatus(status int32) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["status"] = status
}

func (s *SensorStateArg) Status() int32 {
	if ret, ok := s.args["status"].(int32); ok {
		return ret
	} else {
		return 0
	}
}
