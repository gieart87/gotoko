package utils

import "encoding/json"

func PrintJSON(data interface{}) string {
	val, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return ""
	}

	return string(val)
}
