package api

import (
	"github.com/valyala/fasthttp"
	"fmt"
	//"encoding/base64"
	"encoding/json"
	"database/sql"
	"log"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	Login string
	Password string
	ConnectionLimit int
	ConnectionBudget int
}

func (u *user) stringify() string {
	if u.Password == "" {
		return "\"" + u.Login + "\" ,\"\" ," + strconv.Itoa(u.ConnectionLimit) + ", " +
			strconv.Itoa(u.ConnectionBudget)
	}

	return "\"" + u.Login + "\", \"" + u.Password + "\", " + strconv.Itoa(u.ConnectionLimit) + ", " +
		strconv.Itoa(u.ConnectionBudget)
}

func AddUser (cfg *SQLCredStorage) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		db, err := sql.Open(cfg.DBType, cfg.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		var userDescription user
		json.Unmarshal(ctx.PostBody(), &userDescription)
		if userDescription.Password != "" {
			userDescription.Password = encode(userDescription.Password)
		}
		sqlStmt := "insert into " + cfg.UsersTable + " values (" + userDescription.stringify() + ");"
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare(sqlStmt)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec()
		tx.Commit()
	}
}

func RemoveUser(cfg *SQLCredStorage) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		db, err := sql.Open(cfg.DBType, cfg.ConnectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		login := ctx.URI().LastPathSegment()
		sqlStmt := "delete from " + cfg.UsersTable + " where " + cfg.LoginColumn + " = \"" + string(login) +"\""
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare(sqlStmt)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec()
		tx.Commit()
	}
}

func SetConnectionLimit(cfg *SQLCredStorage) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
	}
}

func UpdateUser(cfg *SQLCredStorage) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
	}
}