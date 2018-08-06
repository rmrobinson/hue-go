package hue

// SensorArg is an argument to configure a sensor.
type SensorArg arg

// Reset clears this configuration.
func (s *SensorArg) Reset() {
	s.args = make(map[string]interface{})
}

// Errors exposes any errors encountered when applying the configuration.
func (s *SensorArg) Errors() map[string]ResponseError {
	return s.errors
}

// SetName saves the specified value to be applied.
func (s *SensorArg) SetName(name string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["name"] = name
}

// Name returns the name option, if set.
func (s *SensorArg) Name() string {
	if ret, ok := s.args["name"].(string); ok {
		return ret
	}
	return ""
}

// SensorConfigArg represents a set of configuration changes to be applied to a sensor.
type SensorConfigArg arg

// Reset clears this configuration.
func (s *SensorConfigArg) Reset() {
	s.args = make(map[string]interface{})
}

// Errors exposes any errors encountered when applying the configuration.
func (s *SensorConfigArg) Errors() map[string]ResponseError {
	return s.errors
}

// SetIsOn sets the sensor to be on.
func (s *SensorConfigArg) SetIsOn(isOn bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["on"] = isOn
}

// IsOn returns the name option, if set.
func (s *SensorConfigArg) IsOn() bool {
	if ret, ok := s.args["on"].(bool); ok {
		return ret
	}
	return false
}

// SetIsReachable sets the specified option to be applied.
func (s *SensorConfigArg) SetIsReachable(isReachable bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["reachable"] = isReachable
}

// IsReachable returns the is reachable option, if set.
func (s *SensorConfigArg) IsReachable() bool {
	if ret, ok := s.args["isReachable"].(bool); ok {
		return ret
	}
	return false
}

// SetBatteryLevel sets the specified option to be applied.
func (s *SensorConfigArg) SetBatteryLevel(level uint8) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["battery"] = level
}

// BatteryLevel returns the specified option, if set.
func (s *SensorConfigArg) BatteryLevel() uint8 {
	if ret, ok := s.args["battery"].(uint8); ok {
		return ret
	}
	return 0
}

// SetAlert sets the specified option to be applied.
func (s *SensorConfigArg) SetAlert(alert string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["alert"] = alert
}

// Alert returns the specified option, if set.
func (s *SensorConfigArg) Alert() string {
	if ret, ok := s.args["alert"].(string); ok {
		return ret
	}
	return ""
}

// SetURL sets the specified option to be applied.
func (s *SensorConfigArg) SetURL(url string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["url"] = url
}

// URL returns the specified option, if set.
func (s *SensorConfigArg) URL() string {
	if ret, ok := s.args["url"].(string); ok {
		return ret
	}
	return ""
}

// SetLatitude sets the specified option to be applied.
func (s *SensorConfigArg) SetLatitude(latitude string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["lat"] = latitude
}

// Latitude returns the specified option, if set.
func (s *SensorConfigArg) Latitude() string {
	if ret, ok := s.args["lat"].(string); ok {
		return ret
	}
	return ""
}

// SetLongitude sets the specified option to be applied.
func (s *SensorConfigArg) SetLongitude(longitude string) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["long"] = longitude
}

// Longitude returns the specified option, if set.
func (s *SensorConfigArg) Longitude() string {
	if ret, ok := s.args["long"].(string); ok {
		return ret
	}
	return ""
}

// SetSunriseOffset sets the specified option to be applied.
func (s *SensorConfigArg) SetSunriseOffset(offset int8) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["sunriseoffset"] = offset
}

// SunriseOffset returns the specified option, if set.
func (s *SensorConfigArg) SunriseOffset() int8 {
	if ret, ok := s.args["sunriseoffset"].(int8); ok {
		return ret
	}
	return 0
}

// SetSunsetOffset sets the specified option to be applied.
func (s *SensorConfigArg) SetSunsetOffset(offset int8) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["sunsetoffset"] = offset
}

// SunsetOffset returns the specified option, if set.
func (s *SensorConfigArg) SunsetOffset() int8 {
	if ret, ok := s.args["sunsetoffset"].(int8); ok {
		return ret
	}
	return 0
}

// SetDarkThreshold sets the specified option to be applied.
func (s *SensorConfigArg) SetDarkThreshold(threshold uint16) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["tholddark"] = threshold
}

// DarkThreshold returns the specified option, if set.
func (s *SensorConfigArg) DarkThreshold() uint16 {
	if ret, ok := s.args["tholddark"].(uint16); ok {
		return ret
	}
	return 0
}

// SetDarkThresholdOffset sets the specified option to be applied.
func (s *SensorConfigArg) SetDarkThresholdOffset(offset uint16) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["tholdoffset"] = offset
}

// DarkThresholdOffset returns the specified option, if set.
func (s *SensorConfigArg) DarkThresholdOffset() uint16 {
	if ret, ok := s.args["tholdoffset"].(uint16); ok {
		return ret
	}
	return 0
}

// SensorStateArg is a set of options to configure a sensor state.
type SensorStateArg arg

// Reset clears any set configuration options,
func (s *SensorStateArg) Reset() {
	s.args = make(map[string]interface{})
}

// Errors exposes any errors encountered when applying the configuration.
func (s *SensorStateArg) Errors() map[string]ResponseError {
	return s.errors
}

// SetIsOpen sets the specified option to be applied.
func (s *SensorStateArg) SetIsOpen(isOpen bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["open"] = isOpen
}

// IsOpen returns the specified option, if set.
func (s *SensorStateArg) IsOpen() bool {
	if ret, ok := s.args["open"].(bool); ok {
		return ret
	}
	return false
}

// SetIsPresent sets the specified option to be applied.
func (s *SensorStateArg) SetIsPresent(isPresent bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["presence"] = isPresent
}

// IsPresent returns the specified option, if set.
func (s *SensorStateArg) IsPresent() bool {
	if ret, ok := s.args["presence"].(bool); ok {
		return ret
	}
	return false
}

// SetTemperature sets the specified option to be applied.
func (s *SensorStateArg) SetTemperature(temperature int32) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["temperature"] = temperature
}

// Temperature returns the specified option, if set.
func (s *SensorStateArg) Temperature() int32 {
	if ret, ok := s.args["temperature"].(int32); ok {
		return ret
	}
	return 0
}

// SetHumidity sets the specified option to be applied.
func (s *SensorStateArg) SetHumidity(humidity int32) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["humidity"] = humidity
}

// Humidity returns the specified option, if set.
func (s *SensorStateArg) Humidity() int32 {
	if ret, ok := s.args["humidity"].(int32); ok {
		return ret
	}
	return 0
}

// SetLightLevel sets the specified option to be applied.
func (s *SensorStateArg) SetLightLevel(lightLevel uint16) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["lightlevel"] = lightLevel
}

// LightLevel returns the specified option, if set.
func (s *SensorStateArg) LightLevel() uint16 {
	if ret, ok := s.args["lightlevel"].(uint16); ok {
		return ret
	}
	return 0
}

// SetFlag sets the specified option to be applied.
func (s *SensorStateArg) SetFlag(flag bool) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["flag"] = flag
}

// Flag returns the specified option, if set.
func (s *SensorStateArg) Flag() bool {
	if ret, ok := s.args["flag"].(bool); ok {
		return ret
	}
	return false
}

// SetStatus sets the specified option to be applied.
func (s *SensorStateArg) SetStatus(status int32) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}

	s.args["status"] = status
}

// Status returns the specified option, if set.
func (s *SensorStateArg) Status() int32 {
	if ret, ok := s.args["status"].(int32); ok {
		return ret
	}
	return 0
}
