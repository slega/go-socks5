package api

import (
	"github.com/valyala/fasthttp"
	"fmt"
	//"encoding/base64"
	"encoding/json"
	"database/sql"
	"log"
	"strconv"
)

type user struct {
	Login string
	Password string
	Connection_Limit int
	Connection_Budget int
}

func (u *user) stringify() string {
	if u.Password == "" {
		return "\"" + u.Login + "\" ,\"\" ," + strconv.Itoa(u.Connection_Limit) + ", " +
			strconv.Itoa(u.Connection_Budget)
	}

	return "\"" + u.Login + "\", \"" + u.Password + "\", " + strconv.Itoa(u.Connection_Limit) + ", " +
		strconv.Itoa(u.Connection_Budget)
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
		sqlStmt := "insert into users values (" + userDescription.stringify() + ");"
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
		sqlStmt := "delete from users where login = \"" + string(login) +"\""
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