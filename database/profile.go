package database

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type DbProfile struct {
	Id             int
	Username       string
	MembershipType int
	MembershipId   string
	Json           string
}

func GetProfilesByUsername(db *sql.DB, username string) ([]*DbProfile, error) {
	log.Printf("Getting Profile data from DB for: %v\n", username)

	statement, err := db.Prepare("SELECT * FROM profile WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("GetProfilesByUsername: could not create query: %w", err)
	}

	rows, err := statement.Query(username)
	if err != nil {
		return nil, fmt.Errorf("GetProfilesByUsername: could not execute query: %w", err)
	}

	defer rows.Close()

	profiles, err := getProfilesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetProfilesByUsername: could not handle query result: %w", err)
	}

	return profiles, nil
}

func GetProfileByMembershipId(db *sql.DB, membershipId string) ([]*DbProfile, error) {
	statement, err := db.Prepare("SELECT * FROM profile WHERE membership_id = ?")
	if err != nil {
		return nil, fmt.Errorf("GetProfileByMembershipId: could not create query: %w", err)
	}

	rows, err := statement.Query(membershipId)
	if err != nil {
		return nil, fmt.Errorf("GetProfileByMembershipId: could not create profile table: %w", err)
	}

	defer rows.Close()

	profiles, err := getProfilesFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("GetProfileByMembershipId: could not handle query result: %w", err)
	}

	return profiles, nil
}

func CreateProfile(db *sql.DB, username string, membershipType int, membershipId string, json string) error {
	statement, err := db.Prepare("INSERT INTO profile(username, membership_type, membership_id, json) values(?,?,?,?)")
	if err != nil {
		return fmt.Errorf("CreateProfile: could not create query: %w", err)
	}

	_, err = statement.Exec(username, membershipType, membershipId, json)
	if err != nil {
		return fmt.Errorf("CreateProfile: could not execute insert: %w", err)
	}

	return nil
}

func UpdateProfileById(db *sql.DB, id int, username string, membershipType int, membershipId string, json string) error {
	statement, err := db.Prepare("UPDATE profile SET username = ?, membership_type = ?, membership_id = ?, json = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("UpdateProfileById: could not create query: %w", err)
	}

	_, err = statement.Exec(username, membershipType, membershipId, json, id)
	if err != nil {
		return fmt.Errorf("UpdateProfileById: could not execute update: %w", err)
	}

	return nil
}

func getProfilesFromDbRows(rows *sql.Rows) ([]*DbProfile, error) {
	var profiles []*DbProfile
	var profileId int
	var profileUsername string
	var profileMembershipType int
	var profileMembershipId string
	var profileJson string

	for rows.Next() {
		err := rows.Scan(&profileId, &profileUsername, &profileMembershipType, &profileMembershipId, &profileJson)
		if err != nil {
			return nil, fmt.Errorf("getProfilesFromDbRows: could not scan row: %w", err)
		}

		profile := &DbProfile{}
		profile.Id = profileId
		profile.Username = profileUsername
		profile.MembershipType = profileMembershipType
		profile.MembershipId = profileMembershipId
		profile.Json = profileJson

		profiles = append(profiles, profile)
	}

	return profiles, nil
}
