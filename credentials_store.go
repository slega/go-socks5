package socks5

iimport (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
) 

type SQLCredentials struct {
	table            string
	dbType string
	connectionString string
	credentials      struct {
		login    string
		password string
	}
}

func (s *SQLCredentials) Valid(userLogin, userPassword string) bool {
	db, err := sql.Open(s.dbType, s.connectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sqlStmt := "select * from users where (user like " + userLogin + ")";
	rows, err = db.Query(sqlStmt)

	defer rows.Close()

	for rows.Next() {
		var registeredLogin, registeredPassword string
		err = rows.Scan(&registeredLogin, &registeredPassword)
		if err != nil {
			log.Fatal(err)
		}
		if registeredPassword != userPassword {
			return false
		}
	}

	return true
}
