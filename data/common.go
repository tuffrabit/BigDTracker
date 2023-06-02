package data

import (
	"encoding/json"
	"errors"
	"fmt"
)

func stripApiRepsonseJson(responseJson string) (string, error) {
	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(responseJson), &dat); err != nil {
		return "", fmt.Errorf("stripApiRepsonseJson: could not unmarshal json: %w", err)
	}

	if responseData, exists := dat["Response"]; exists {
		newJson, err := json.Marshal(responseData)
		if err != nil {
			return "", fmt.Errorf("stripApiRepsonseJson: could not marshal new json: %w", err)
		}

		return string(newJson), nil
	} else {
		return "", errors.New("stripApiRepsonseJson: json is malformed")
	}
}
