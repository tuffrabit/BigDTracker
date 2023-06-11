package activities

type Activities struct {
	Response RepsonseBody `json:"Response"`
	rawJson  string
}

type RepsonseBody struct {
	Activities []RepsonseActivity `json:"activities"`
}

type RepsonseActivity struct {
	Details RepsonseDetails `json:"activityDetails"`
}

type RepsonseDetails struct {
	ReferenceId int    `json:"referenceId"`
	InstanceId  string `json:"instanceId"`
}

func (activities *Activities) SetRawJson(rawJson string) {
	activities.rawJson = rawJson
}

func (activities *Activities) GetRawJson() string {
	return activities.rawJson
}
