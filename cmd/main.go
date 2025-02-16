package main

import (
	httphandler "avito2015/api/v1"
	"avito2015/internal/db"
	"avito2015/internal/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	err := db.Init()
	if err != nil {
		panic("db init error")
	}
	defer db.DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/auth", httphandler.AuthHandler).Methods("POST")
	r.HandleFunc("/api/sendCoin", middleware.Auth(httphandler.SendCoinHandler)).Methods("POST")
	r.HandleFunc("/api/info", middleware.Auth(httphandler.InfoHandler)).Methods("GET")
	r.HandleFunc("/api/buy/{item}", middleware.Auth(httphandler.BuyItemHandler)).Methods("GET")

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
