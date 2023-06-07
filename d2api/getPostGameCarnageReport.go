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

func (api *Api) GetPostGameCarnageReport(instanceId string) (*PostGameCarnageReportRepsonse, error) {
	log.Printf("Getting PostGameCarnageReport data from Bungie for: %v\n", instanceId)

	responseBody, err := api.DoGetRequest(fmt.Sprintf("Stats/PostGameCarnageReport/%v/", instanceId))
	if err != nil {
		return nil, fmt.Errorf("Api.GetPostGameCarnageReport: could not create bungie api request: %w", err)
	}

	jsonResponse := &PostGameCarnageReportRepsonse{}
	err = json.Unmarshal(*responseBody, jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("Api.GetPostGameCarnageReport: could not json decode response: %w", err)
	}

	jsonResponse.Json = string(*responseBody)

	return jsonResponse, nil
}
