package main

import (
	"net/http"
	"github.com/coralproject/pillar/server/log"
	"github.com/coralproject/pillar/server/web"
)

func main() {

	router := web.NewRouter()

	log.Logger.Print(http.ListenAndServe(":8080", router))
}