package d2api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func GetPostGameCarnageReport(apiKey string, instanceId string) (*PostGameCarnageReportRepsonse, error) {
	log.Printf("Getting PostGameCarnageReport data from Bungie for: %v\n", instanceId)

	request, err := http.NewRequest(
		"GET",
		baseUrl+"Stats/PostGameCarnageReport/"+instanceId+"/",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReport: could not create bungie api request: %w", err)
	}

	request.Header.Add("X-API-KEY", apiKey)

	client := &http.Client{}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReport: request to bungie api failed: %w", err)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetPostGameCarnageReport: response from bungie api returned http %v", httpResponse.StatusCode)
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReport: could not read response body: %w", err)
	}

	jsonResponse := &PostGameCarnageReportRepsonse{}
	err = json.Unmarshal(body, jsonResponse)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReport: could not json decode response: %w", err)
	}

	jsonResponse.Json = string(body)

	return jsonResponse, nil
}
