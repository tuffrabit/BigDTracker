package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	Db *sql.DB
}

type DbHandler interface {
	Init(db *sql.DB) error
}

type DbHandlers struct {
	Profile               *DbProfile
	Activity              *DbActivity
	ActivityHistory       *DbActivityHistory
	PostGameCarnageReport *DbPostGameCarnageReport
}

func (db *Db) GetDbHandlers() (*DbHandlers, error) {
	initDb := false

	if _, err := os.Stat("./data.db"); err != nil {
		initDb = true
	}

	dsn := "./data.db"
	sqlDb, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("Db.GetDbHandlers: could not open data.db: %w", err)
	}

	db.Db = sqlDb

	if initDb {
		if err = db.initializeDb(); err != nil {
			return nil, fmt.Errorf("Db.GetDbHandlers: could not initialize data.db: %w", err)
		}
	}

	dbProfile := &DbProfile{}
	err = dbProfile.Init(db.Db)
	if err != nil {
		return nil, fmt.Errorf("Db.GetDbHandlers: could not initialize profile handler: %w", err)
	}

	dbActivity := &DbActivity{}
	err = dbActivity.Init(db.Db)
	if err != nil {
		return nil, fmt.Errorf("Db.GetDbHandlers: could not initialize activity handler: %w", err)
	}

	dbActivityHistory := &DbActivityHistory{}
	err = dbActivityHistory.Init(db.Db)
	if err != nil {
		return nil, fmt.Errorf("Db.GetDbHandlers: could not initialize activityhistory handler: %w", err)
	}

	dbPostGameCarnageReport := &DbPostGameCarnageReport{}
	err = dbPostGameCarnageReport.Init(db.Db)
	if err != nil {
		return nil, fmt.Errorf("Db.GetDbHandlers: could not initialize post game carnage report handler: %w", err)
	}

	dbHandlers := &DbHandlers{}
	dbHandlers.Profile = dbProfile
	dbHandlers.Activity = dbActivity
	dbHandlers.ActivityHistory = dbActivityHistory
	dbHandlers.PostGameCarnageReport = dbPostGameCarnageReport

	return dbHandlers, nil
}

func (db *Db) initializeDb() error {
	createStatements, err := db.getInitializeDbTableCreateStatements()
	if err != nil {
		return fmt.Errorf("Db.InitializeDb: could not get table create statements: %w", err)
	}

	for _, statement := range createStatements {
		_, err := statement.Exec()
		if err != nil {
			return fmt.Errorf("Db.InitializeDb: could not execute table create statement: %w", err)
		}
	}

	createStatements, err = db.getInitializeDbIndexCreateStatements()
	if err != nil {
		return fmt.Errorf("Db.InitializeDb: could not get index create statements: %w", err)
	}

	for _, statement := range createStatements {
		_, err := statement.Exec()
		if err != nil {
			return fmt.Errorf("Db.InitializeDb: could not execute index create statement: %w", err)
		}
	}

	return nil
}

func (db *Db) getInitializeDbTableCreateStatements() ([]*sql.Stmt, error) {
	var statements []*sql.Stmt

	statement, err := db.Db.Prepare(`
		CREATE TABLE profile (
			id INTEGER PRIMARY KEY,
			username TEXT NOT NULL,
			membership_type INTEGER NOT NULL,
			membership_id TEXT NOT NULL UNIQUE,
			json TEXT NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("Db.getInitializeDbTableCreationStatements: could not create profile table query: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Db.Prepare(`
		CREATE TABLE activity (
			id INTEGER PRIMARY KEY,
			instance_id TEXT NOT NULL,
			membership_ids TEXT NOT NULL,
			membership_type INTEGER NOT NULL,
			character_ids TEXT NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("Db.getInitializeDbTableCreationStatements: could not create activity table: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Db.Prepare(`
		CREATE TABLE activity_history (
			id INTEGER PRIMARY KEY,
			membership_id TEXT NOT NULL,
			membership_type INTEGER NOT NULL,
			character_id TEXT NOT NULL,
			activity_count INTEGER NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("Db.getInitializeDbTableCreationStatements: could not create activity_history table: %w", err)
	}

	statements = append(statements, statement)

	statement, err = db.Db.Prepare(`
		CREATE TABLE post_game_carnage_report (
			id INTEGER PRIMARY KEY,
			instance_id TEXT NOT NULL UNIQUE,
			json TEXT NOT NULL
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("Db.getInitializeDbTableCreationStatements: could not create post_game_carnage_report table: %w", err)
	}

	statements = append(statements, statement)

	return statements, nil
}

func (db *Db) getInitializeDbIndexCreateStatements() ([]*sql.Stmt, error) {
	var statements []*sql.Stmt

	statement, err := db.Db.Prepare(`
			CREATE UNIQUE INDEX activity_history_unique_index
			ON activity_history (membership_id, membership_type, character_id
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("Db.getInitializeDbIndexCreateStatements: could not create activity_history_unique_index index: %w", err)
	}

	statements = append(statements, statement)

	return statements, nil
}
