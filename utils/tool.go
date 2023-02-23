package utils

import "encoding/json"

func SortJsonStr(str string) string {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(str), &jsonMap); err != nil {
		return str
	}
	data, err := json.Marshal(str)
	if err != nil {
		return str
	}
	return string(data)
}
