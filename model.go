package main

import (
	"database/sql"
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

func (s *secret) getsecret() error {
	return DB.QueryRow("SELECT `secretID`, `expiration`, `contents`, `ContentsMeta`, `userID`, `name` FROM `secrets` WHERE id=?",
		s.SecretID).Scan(&s.SecretID, &s.Expiration, &s.Contents, &s.ContentsMeta, &s.UserID, &s.Name)
}

func (s *secret) updatesecret() error {
	_, err := DB.Exec("UPDATE `secrets` SET `secretID`=?,`expiration`=?,`contents`=?,`ContentsMeta`=?,`userID`=?,`name`=? WHERE `secretID`=?", &s.SecretID, &s.Expiration, &s.Contents, &s.ContentsMeta, &s.UserID, &s.Name)
	return err
}

func (s *secret) deletesecret() error {
	_, err := DB.Exec("DELETE FROM secrets WHERE secretID=?", s.SecretID)
	return err
}

func (s *secret) createsecret() error {

	_, err := DB.Exec("INSERT INTO `secrets`(`secretID`, `expiration`, `contents`, `ContentsMeta`, `userID`, `name`) VALUES (?,?,?,?,?,?)", &s.SecretID, &s.Expiration, &s.Contents, &s.ContentsMeta, &s.UserID, &s.Name)
	return err
}
