package postgamecarnagereport

type PostGameCarnageReport struct {
	Response RepsonseBody `json:"Response"`
	rawJson  string
}

type RepsonseBody struct {
	Entries []RepsonseEntry `json:"entries"`
}

type RepsonseEntry struct {
	CharacterId string           `json:"characterId"`
	Extended    RepsonseExtended `json:"extended"`
}

type RepsonseExtended struct {
	Weapons []RepsonseWeapon `json:"weapons"`
}

type RepsonseWeapon struct {
	ReferenceId int                  `json:"referenceId"`
	Values      RepsonseWeaponValues `json:"values"`
}

type RepsonseWeaponValues struct {
	UniqueWeaponKills          RepsonseBasicValueContainer `json:"uniqueWeaponKills"`
	UniqueWeaponPrecisionKills RepsonseBasicValueContainer `json:"uniqueWeaponPrecisionKills"`
}

type RepsonseBasicValueContainer struct {
	Basic RepsonseBasicValue `json:"basic"`
}

type RepsonseBasicValue struct {
	Value float64 `json:"value"`
}

func (postGameCarnageReport *PostGameCarnageReport) SetRawJson(rawJson string) {
	postGameCarnageReport.rawJson = rawJson
}

func (postGameCarnageReport *PostGameCarnageReport) GetRawJson() string {
	return postGameCarnageReport.rawJson
}
