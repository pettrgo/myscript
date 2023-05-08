package utils

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"os"
)

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

func UnsafeMarshal(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(data)
}

func NewCsvFile(ctx context.Context, filePath string) (*csv.Writer, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	return csv.NewWriter(file), nil
}
