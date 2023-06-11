package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

type DbPostGameCarnageReportData struct {
	Id         int
	InstanceId string
	Json       string
}

type DbPostGameCarnageReport struct {
	Db
	PreparedStatements *DbPostGameCarnageReportPreparedStatements
}

type DbPostGameCarnageReportPreparedStatements struct {
	GetByInstanceId *sql.Stmt
	Create          *sql.Stmt
	UpdateById      *sql.Stmt
}

func (dbPostGameCarnageReport *DbPostGameCarnageReport) Init(db *sql.DB) error {
	preparedStatements := &DbPostGameCarnageReportPreparedStatements{}
	statement, err := db.Prepare("SELECT * FROM post_game_carnage_report WHERE instance_id = ?")
	if err != nil {
		return fmt.Errorf("DbPostGameCarnageReport.Init: could not create query: %w", err)
	}

	preparedStatements.GetByInstanceId = statement

	statement, err = db.Prepare("INSERT INTO post_game_carnage_report(instance_id, json) values(?,?)")
	if err != nil {
		return fmt.Errorf("DbPostGameCarnageReport.Init: could not create query: %w", err)
	}

	preparedStatements.Create = statement

	statement, err = db.Prepare("UPDATE post_game_carnage_report SET instance_id = ?, json = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("DbPostGameCarnageReport.Init: could not create query: %w", err)
	}

	preparedStatements.UpdateById = statement
	dbPostGameCarnageReport.PreparedStatements = preparedStatements
	dbPostGameCarnageReport.Db.Db = db

	return nil
}

func (dbPostGameCarnageReport *DbPostGameCarnageReport) GetPostGameCarnageReportsByInstanceId(instanceId string) ([]*DbPostGameCarnageReportData, error) {
	log.Printf("Getting PostGameCarnageReport data from DB for: %v\n", instanceId)

	rows, err := dbPostGameCarnageReport.PreparedStatements.GetByInstanceId.Query(instanceId)
	if err != nil {
		return nil, fmt.Errorf("GetPostGameCarnageReportByInstanceId: could not create post_game_carnage_report query: %w", err)
	}

	defer rows.Close()

	postGameCarnageReports, err := dbPostGameCarnageReport.getPostGameCarnageReportFromDbRows(rows)
	if err != nil {
		return nil, fmt.Errorf("DbPostGameCarnageReport.GetPostGameCarnageReportByInstanceId: could not handle query result: %w", err)
	}

	return postGameCarnageReports, nil
}

func (dbPostGameCarnageReport *DbPostGameCarnageReport) CreatePostGameCarnageReport(instanceId string, json string) error {
	log.Printf("Inserting PostGameCarnageReport data into DB for: %v\n", instanceId)

	_, err := dbPostGameCarnageReport.PreparedStatements.Create.Exec(instanceId, json)
	if err != nil {
		return fmt.Errorf("DbPostGameCarnageReport.CreatePostGameCarnageReport: could not execute insert: %w", err)
	}

	return nil
}

func (dbPostGameCarnageReport *DbPostGameCarnageReport) UpdatePostGameCarnageReportById(id int, instanceId string, json string) error {
	log.Printf("Updating PostGameCarnageReport data into DB for: %v\n", instanceId)

	_, err := dbPostGameCarnageReport.PreparedStatements.UpdateById.Exec(instanceId, json, id)
	if err != nil {
		return fmt.Errorf("DbPostGameCarnageReport.UpdatePostGameCarnageReportById: could not execute update: %w", err)
	}

	return nil
}

func (dbPostGameCarnageReport *DbPostGameCarnageReport) getPostGameCarnageReportFromDbRows(rows *sql.Rows) ([]*DbPostGameCarnageReportData, error) {
	var postGameCarnageReports []*DbPostGameCarnageReportData
	var id int
	var instanceId string
	var json string

	for rows.Next() {
		err := rows.Scan(&id, &instanceId, &json)
		if err != nil {
			return nil, fmt.Errorf("DbPostGameCarnageReport.getPostGameCarnageReportFromDbRows: could not scan row: %w", err)
		}

		postGameCarnageReport := &DbPostGameCarnageReportData{}
		postGameCarnageReport.Id = id
		postGameCarnageReport.InstanceId = instanceId
		postGameCarnageReport.Json = json

		postGameCarnageReports = append(postGameCarnageReports, postGameCarnageReport)
	}

	return postGameCarnageReports, nil
}
