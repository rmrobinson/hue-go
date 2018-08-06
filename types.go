package hue

import "encoding/json"

type arg struct {
	args    map[string]interface{}
	errors  map[string]ResponseError
	success map[string]interface{}
}

var responseErrorTypes = map[int]string{
	1:   "Unauthorized User",
	2:   "Body contains invalid JSON",
	3:   "Resource not available",
	4:   "Method not available for resource",
	5:   "Missing parameter in body",
	6:   "Parameter not available",
	7:   "Invalid value for parameter",
	8:   "Parameter not modifiable",
	11:  "Too many items in list",
	12:  "Portal connection required",
	901: "Internal error",
}

// ResponseError is the error message returned if the given entry is invalid.
// The address refers to the component which failed to change; the Type maps to errorTypes above.
type ResponseError struct {
	Type        int    `json:"type"`
	Description string `json:"description"`
	Address     string `json:"address"`
}

// One of the entries in the response array returned by the API to a PUT/POST request.
type responseEntry struct {
	Success map[string]*json.RawMessage `json:"success"`
	Error   ResponseError               `json:"error"`
}

// The overall response array returned by the API to a PUT/POST request.
type responseEntries []json.RawMessage
