package api

import (
	"database/sql"
	"log"
	"encoding/base64"
	_ "github.com/mattn/go-sqlite3"
)

type SQLCredStorage struct {
	UsersTable       string
	LoginColumn      string
	DBType           string
	ConnectionString string
}

func encode(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func (s *SQLCredStorage) Valid(user, password string) bool {
	db, err := sql.Open(s.DBType, s.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sqlStmt := "select * from " + s.UsersTable + " where (login like \"" + user + "\")"
	rows, err := db.Query(sqlStmt)

	defer rows.Close()

	for rows.Next() {
		var registeredLogin, registeredPassword string
		var limit, budget int
		err = rows.Scan(&registeredLogin, &registeredPassword, &limit, &budget)
		if err != nil {
			log.Fatal(err)
		}

		if registeredPassword != encode(password) {
			return false
		} else {
			return true
		}
	}

	return false
}
