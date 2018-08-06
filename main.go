package main

import (
	"sync"
	"github.com/buaazp/fasthttprouter"
	"log"
	"github.com/valyala/fasthttp"
	"go-socks5/socks5"
	"go-socks5/api"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		cfg := api.SQLCredStorage {
			"users",
			"sqlite3",
			"/Users/slega/go/src/go-socks5/proxy_users.db",
		}

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
		conf.Credentials = &api.SQLCredStorage{}
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
