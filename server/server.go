package main

import (
	"net/http"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"

	"github.com/coralproject/pillar/server/pkg/stats"
	"github.com/coralproject/pillar/server/web"
)

func init() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)
}

func main() {

	log.Dev("startup", "main", "Start")

	stats.Init()

	router := web.NewRouter()

	log.Error("startup", "main", http.ListenAndServe(":8080", router), "Listening")
}
