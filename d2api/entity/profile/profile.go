package profile

type Profile struct {
	Response RepsonseBody `json:"Response"`
	rawJson  string
}

type RepsonseBody struct {
	Profile RepsonseProfile `json:"profile"`
}

type RepsonseProfile struct {
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	CharacterIds []string `json:"characterIds"`
}

func (profile *Profile) SetRawJson(rawJson string) {
	profile.rawJson = rawJson
}

func (profile *Profile) GetRawJson() string {
	return profile.rawJson
}
