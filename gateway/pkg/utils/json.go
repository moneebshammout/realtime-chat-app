package utils

import (
	"encoding/json"
	"os"
)

func ParseJsonFile(filePath string, dto interface{}) (interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return dto, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(dto)
	if err != nil {
		return nil, err
	}

	return dto, nil
}
