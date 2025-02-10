package database

import (
	"database/sql"
	"log"
	"time"
)

type Query struct {
	db   *sql.DB
	Time *time.Location
}

func NewQuery(db *sql.DB) *Query {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}

	return &Query{
		db:   db,
		Time: loc,
	}
}

func (q *Query) CreateTables() error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS customername (
			customer_name VARCHAR(255) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS bsno (
			bs_number VARCHAR(50) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS partcode (
			part_code VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS unit (	
			unit_name VARCHAR(100) NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS postatusdd (
			status VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS concernsonorder (
			concern VARCHAR(100) NOT NULL
		)`,
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			log.Printf("Failed to execute query: %s", query)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	log.Println("All tables created successfully.")
	return nil
}
