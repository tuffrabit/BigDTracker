package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tuffrabit/BigDTracker/d2api"
	"github.com/tuffrabit/BigDTracker/data"
	"github.com/tuffrabit/BigDTracker/database"
)

func main() {
	mainArgs, err := validateStartup()
	if err != nil {
		panic(err)
	}

	apiKey, err := getApiKey()
	if err != nil {
		panic(err)
	}

	log.SetOutput(io.Discard)

	if mainArgs.LogLocation == 1 {
		log.SetOutput(os.Stdout)
	}

	api := &d2api.Api{}
	api.Init(apiKey)

	db := &database.Db{}
	dbHandlers, err := db.GetDbHandlers()
	if err != nil {
		panic(err)
	}

	profiles, err := data.GetProfilesByUsername(dbHandlers.Profile, api, mainArgs.Username)
	if err != nil {
		panic(err)
	}

	bigDKillCount := 0.0

	for _, profile := range profiles {
		for _, characterId := range profile.Profile.Data.CharacterIds {
			activities, err := data.GetActivities(
				dbHandlers,
				api,
				profile.DBProfile.MembershipId,
				profile.DBProfile.MembershipType,
				characterId,
			)
			if err != nil {
				panic(err)
			}

			for _, activity := range activities {
				log.Printf("Processing Activity #: %v\n", activity.InstanceId)

				postGameCarnageReports, err := data.GetPostGameCarnageReportsByInstanceId(dbHandlers.PostGameCarnageReport, api, activity.InstanceId)
				if err != nil {
					log.Println(fmt.Errorf("main: could not get post game carnage report for %v: %w", activity.InstanceId, err))
					break
				}

				for _, postGameCarnageReport := range postGameCarnageReports {
					for _, entry := range postGameCarnageReport.Entries {
						if entry.CharacterId == characterId {
							for _, weapon := range entry.Extended.Weapons {
								if weapon.ReferenceId == d2api.BigDApiHash {
									count := weapon.Values.UniqueWeaponKills.Basic.Value + weapon.Values.UniqueWeaponPrecisionKills.Basic.Value
									bigDKillCount = bigDKillCount + count
								}
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("Big D PvP kills: %v", bigDKillCount)
}

type MainArgs struct {
	Username    string
	LogLocation int
}

func validateStartup() (*MainArgs, error) {
	var user string

	if len(os.Args) > 1 {
		user = os.Args[1]

		if user == "" {
			return nil, errors.New("validateStartup: you must supply a bungie username")
		}

		if !strings.Contains(user, "#") {
			return nil, errors.New("validateStartup: the supplied bungie username is an unvalid format")
		}
	} else {
		return nil, errors.New("validateStartup: you must supply a bungie username")
	}

	var logLocation int

	if len(os.Args) > 2 {
		var err error
		logLocation, err = strconv.Atoi(os.Args[2])
		if err != nil {
			return nil, fmt.Errorf("validateStartup: could not parse log location argument: %w", err)
		}
	}

	var mainArgs = &MainArgs{}
	mainArgs.Username = user
	mainArgs.LogLocation = logLocation

	return mainArgs, nil
}

func getApiKey() (string, error) {
	apiKey, exists := os.LookupEnv("BUNGIE_API_KEY")
	if !exists {
		return "", errors.New("getApiKey: bungie api key env variable is not set")
	}

	if apiKey == "" {
		return "", errors.New("getApiKey: bungie api key cannot be blank")
	}

	return apiKey, nil
}
