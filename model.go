package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB ...
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
}

func (p *secret) getsecret(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *secret) updatesecret(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *secret) deletesecret(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *secret) createsecret(db *sql.DB) error {
	return errors.New("Not implemented")
}
