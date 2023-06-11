package d2data

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tuffrabit/BigDTracker/d2api"
	"github.com/tuffrabit/BigDTracker/database"
)

type Data struct {
	DbHandlers *database.DbHandlers
	Api        *d2api.Api
}

func NewData(dbHandlers *database.DbHandlers, api *d2api.Api) *Data {
	return &Data{
		DbHandlers: dbHandlers,
		Api:        api,
	}
}

func (data *Data) Init(dbHandlers *database.DbHandlers, api *d2api.Api) {
	data.DbHandlers = dbHandlers
	data.Api = api
}

func (data *Data) stripApiRepsonseJson(responseJson string) (string, error) {
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
