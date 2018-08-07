package main

import (
	"sync"
	"github.com/buaazp/fasthttprouter"
	"log"
	"github.com/valyala/fasthttp"
	"go-socks5/socks5"
	"go-socks5/api"
	_ "github.com/mattn/go-sqlite3"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-config"
	"flag"
	"fmt"
)

func main() {
	configFlag := flag.String("conf", "", "location of configuration file")
	flag.Parse()
	if *configFlag == "" {
		fmt.Println("No configuration file was specified\n")
		return
	}
	config.Load(file.NewSource(
		file.WithPath(*configFlag),
	))

	var conf struct {
		UsersTable string `json:"usersTable"`
		LoginColumn string `json:"loginColumn"`
		DBtype string `json:"DBtype"`
		ConnectionString string `json:"connectionString"`
	}

	if err := config.Scan(&conf); err != nil {
		panic(err)
	}

	cfg := api.SQLCredStorage {
		UsersTable: conf.UsersTable,
		LoginColumn: conf.LoginColumn,
		DBType: conf.DBtype,
		ConnectionString: conf.ConnectionString,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		router := fasthttprouter.New()
		router.POST("/user/", api.AddUser(&cfg))
		router.DELETE("/user/:name", api.RemoveUser(&cfg))
		router.PUT("/user/:name", api.UpdateUser(&cfg))
		router.PUT("/limit/", api.SetConnectionLimit(&cfg))

		log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
	}()

	go func() {
		defer wg.Done()

		conf := &socks5.Config{}
		conf.Credentials = &cfg

		server, err := socks5.New(conf)
		if err != nil {
			panic(err)
		}

		// Create SOCKS5 proxy on localhost port 8000
		if err := server.ListenAndServe("tcp", "127.0.0.1:8000"); err != nil {
			panic(err)
		}
	}()


	wg.Wait()
	return
}
