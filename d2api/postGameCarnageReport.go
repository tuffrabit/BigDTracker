package d2api

import (
	"encoding/json"
	"fmt"
	"log"
)

type PostGameCarnageReportRepsonse struct {
	Response PostGameCarnageReportRepsonseResponse `json:"Response"`
	Json     string
}

type PostGameCarnageReportRepsonseResponse struct {
	Entries []PostGameCarnageReportRepsonseEntry `json:"entries"`
}

type PostGameCarnageReportRepsonseEntry struct {
	CharacterId string                                `json:"characterId"`
	Extended    PostGameCarnageReportRepsonseExtended `json:"extended"`
}

type PostGameCarnageReportRepsonseExtended struct {
	Weapons []PostGameCarnageReportRepsonseWeapon `json:"weapons"`
}

type PostGameCarnageReportRepsonseWeapon struct {
	ReferenceId int                                       `json:"referenceId"`
	Values      PostGameCarnageReportRepsonseWeaponValues `json:"values"`
}

type PostGameCarnageReportRepsonseWeaponValues struct {
	UniqueWeaponKills          PostGameCarnageReportRepsonseBasicValueContainer `json:"uniqueWeaponKills"`
	UniqueWeaponPrecisionKills PostGameCarnageReportRepsonseBasicValueContainer `json:"uniqueWeaponPrecisionKills"`
}

type PostGameCarnageReportRepsonseBasicValueContainer struct {
	Basic PostGameCarnageReportRepsonseBasicValue `json:"basic"`
}

type PostGameCarnageReportRepsonseBasicValue struct {
	Value float64 `json:"value"`
}

func (getPostGameCarnageReportRepsonse *PostGameCarnageReportRepsonse) UnmarshalHttpResponseBody(responseBody []byte) error {
	err := json.Unmarshal(responseBody, getPostGameCarnageReportRepsonse)
	if err != nil {
		return fmt.Errorf("PostGameCarnageReportRepsonse.UnmarshalHttpResponseBody: could not json decode response: %w", err)
	}

	getPostGameCarnageReportRepsonse.Json = string(responseBody)

	return nil
}

func (api *Api) GetPostGameCarnageReport(instanceId string) (*PostGameCarnageReportRepsonse, error) {
	log.Printf("Getting PostGameCarnageReport data from Bungie for: %v\n", instanceId)

	entity := &PostGameCarnageReportRepsonse{}
	err := api.DoGetRequest(
		fmt.Sprintf("Stats/PostGameCarnageReport/%v/", instanceId),
		entity,
	)
	if err != nil {
		return nil, fmt.Errorf("Api.GetPostGameCarnageReport: could not create bungie api request: %w", err)
	}

	return entity, nil
}
