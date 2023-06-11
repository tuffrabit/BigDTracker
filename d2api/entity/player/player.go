package player

type Player struct {
	Response []RepsonseBody `json:"Response"`
	rawJson  string
}

type RepsonseBody struct {
	MembershipId                string `json:"membershipId"`
	MembershipType              int    `json:"membershipType"`
	BungieGlobalDisplayNameCode int    `json:"bungieGlobalDisplayNameCode"`
}

func (player *Player) SetRawJson(rawJson string) {
	player.rawJson = rawJson
}

func (player *Player) GetRawJson() string {
	return player.rawJson
}
