package main

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	go func() {
		router := fasthttprouter.New()
		router.POST("/user/:name", AddUser)
		router.DELETE("/user/:name", RemoveUser)
		router.PUT("/user/:name", UpdateUser)
		router.PUT("/limit/", SetConnectionLimit)

		log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
		conf := &socks5.Config{}
		server, err := socks5.New(conf)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		// Create SOCKS5 proxy on localhost port 8000
		if err := server.ListenAndServe("tcp", "127.0.0.1:8000"); err != nil {
			panic(err)
		}
	}()

	return
}
