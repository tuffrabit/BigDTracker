package data

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/tuffrabit/BigDTracker/d2api"
	"github.com/tuffrabit/BigDTracker/database"
)

type ProfileContainer struct {
	Profile                            Profile `json:"profile"`
	ResponseMintedTimestamp            string  `json:"responseMintedTimestamp"`
	SecondaryComponentsMintedTimestamp string  `json:"secondaryComponentsMintedTimestamp"`
	DBProfile                          *database.DbProfile
}

type Profile struct {
	Data    ProfileData `json:"data"`
	Privacy int         `json:"privacy"`
}

type ProfileData struct {
	CharacterIds        []string `json:"characterIds"`
	CurrentGuardianRank int      `json:"currentGuardianRank"`
}

func GetProfilesByUsername(db *sql.DB, apiKey string, username string) ([]*ProfileContainer, error) {
	dbProfiles, err := database.GetProfilesByUsername(db, username)
	if err != nil {
		return nil, fmt.Errorf("GetProfilesByUsername: could not get profiles from db: %w", err)
	}

	if len(dbProfiles) == 0 {
		playerResponse, err := d2api.GetPlayer(apiKey, username)
		if err != nil {
			return nil, fmt.Errorf("GetProfilesByUsername: could not get player data from api: %w", err)
		}

		for _, player := range playerResponse.Response {
			profileResponse, err := d2api.GetProfile(apiKey, player.MembershipId, player.MembershipType)
			if err == nil {
				profileJson, err := stripApiRepsonseJson(profileResponse.Json)
				if err != nil {
					return nil, fmt.Errorf("GetProfilesByUsername: could not strip response json: %w", err)
				}

				err = database.CreateProfile(db, username, player.MembershipType, player.MembershipId, profileJson)
				if err != nil {
					return nil, fmt.Errorf("GetProfilesByUsername: could not insert profile data into db: %w", err)
				}

				dbProfiles, err = database.GetProfilesByUsername(db, username)
				if err != nil {
					return nil, fmt.Errorf("GetProfilesByUsername: could not get profiles from db: %w", err)
				}
			}
		}
	}

	var profiles []*ProfileContainer

	for _, dbProfile := range dbProfiles {
		profile := &ProfileContainer{}

		err := json.Unmarshal([]byte(dbProfile.Json), profile)
		if err != nil {
			return nil, fmt.Errorf("GetProfilesByUsername: could not get unmarshal json data from db: %w", err)
		}

		profile.DBProfile = dbProfile
		profiles = append(profiles, profile)
	}

	return profiles, nil
}
