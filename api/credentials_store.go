package api

import (
	"database/sql"
	"log"
	"encoding/base64"
)

type SQLCredStorage struct {
	UsersTable string
	DBType string
	ConnectionString string
}

func encode(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func (s *SQLCredStorage)  Valid(user, password string) bool  {
	db, err := sql.Open(s.DBType, s.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sqlStmt := "select * from users where (user like " + user + ")";
	rows, err := db.Query(sqlStmt)

	defer rows.Close()

	for rows.Next() {
		var registeredLogin, registeredPassword string
		err = rows.Scan(&registeredLogin, &registeredPassword)
		if err != nil {
			log.Fatal(err)
		}
		if registeredPassword != encode(password) {
			return false
		}
	}

	return true
}
