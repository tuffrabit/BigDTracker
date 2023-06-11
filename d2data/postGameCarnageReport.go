package d2data

import (
	"encoding/json"
	"fmt"

	"github.com/tuffrabit/BigDTracker/database"
)

type PostGameCarnageReport struct {
	Period                  string                       `json:"period"`
	Entries                 []PostGameCarnageReportEntry `json:"entries"`
	DbPostGameCarnageReport *database.DbPostGameCarnageReportData
}

type PostGameCarnageReportEntry struct {
	CharacterId string                        `json:"characterId"`
	Extended    PostGameCarnageReportExtended `json:"extended"`
}

type PostGameCarnageReportExtended struct {
	Weapons []PostGameCarnageReportWeapon `json:"weapons"`
}

type PostGameCarnageReportWeapon struct {
	ReferenceId int                               `json:"referenceId"`
	Values      PostGameCarnageReportWeaponValues `json:"values"`
}

type PostGameCarnageReportWeaponValues struct {
	UniqueWeaponKills          PostGameCarnageReportWeaponValue `json:"uniqueWeaponKills"`
	UniqueWeaponPrecisionKills PostGameCarnageReportWeaponValue `json:"uniqueWeaponPrecisionKills"`
}

type PostGameCarnageReportWeaponValue struct {
	Basic PostGameCarnageReportBasicValue `json:"basic"`
}

type PostGameCarnageReportBasicValue struct {
	DisplayValue string  `json:"displayValue"`
	Value        float64 `json:"value"`
}

func (data *Data) GetPostGameCarnageReportsByInstanceId(instanceId string) ([]*PostGameCarnageReport, error) {
	dbPostGameCarnageReports, err := data.DbHandlers.PostGameCarnageReport.GetPostGameCarnageReportsByInstanceId(instanceId)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReportByInstnaceId: could not get post game carnage report from db: %w", err)
	}

	if len(dbPostGameCarnageReports) == 0 {
		postGameCarnageReportResponse, err := data.Api.PostGameCarnageReportHandler.DoGet(instanceId)
		if err != nil {
			return nil, fmt.Errorf("GetPostGameCarnageReportByInstnaceId: could not get post game carnage report data from api: %w", err)
		}

		postGameCarnageReportJson, err := data.stripApiRepsonseJson(postGameCarnageReportResponse.GetRawJson())
		if err != nil {
			return nil, fmt.Errorf("GetPostGameCarnageReportByInstnaceId: could not strip response json: %w", err)
		}

		err = data.DbHandlers.PostGameCarnageReport.CreatePostGameCarnageReport(instanceId, postGameCarnageReportJson)
		if err != nil {
			return nil, fmt.Errorf("GetPostGameCarnageReportByInstnaceId: could not insert post game carnage report data into db: %w", err)
		}

		dbPostGameCarnageReports, err = data.DbHandlers.PostGameCarnageReport.GetPostGameCarnageReportsByInstanceId(instanceId)
		if err != nil {
			return nil, fmt.Errorf("GetPostGameCarnageReportByInstnaceId: could not get post game carnage report from db: %w", err)
		}
	}

	var postGameCarnageReports []*PostGameCarnageReport

	for _, dbPostGameCarnageReport := range dbPostGameCarnageReports {
		postGameCarnageReport := &PostGameCarnageReport{}

		err := json.Unmarshal([]byte(dbPostGameCarnageReport.Json), postGameCarnageReport)
		if err != nil {
			return nil, fmt.Errorf("GetPostGameCarnageReportsByInstanceId: could not get unmarshal json data from db: %w", err)
		}

		postGameCarnageReport.DbPostGameCarnageReport = dbPostGameCarnageReport
		postGameCarnageReports = append(postGameCarnageReports, postGameCarnageReport)
	}

	return postGameCarnageReports, nil
}
