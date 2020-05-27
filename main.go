package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"nawp.com/util/config"
	"nawp.com/util/database"

	"nawp.com/app"
)

func main() {
	cfg := config.LoadConfig("config.json")
	db := database.InitDB(cfg)

	err := http.ListenAndServe(":"+cfg.PORT, app.Router(mux.NewRouter(), db))

	if err != nil {
		log.Fatal("Unable to start the web server!")
	}
}
