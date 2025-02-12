package main

import (
	httphandler "avito2015/api/v1"
	"avito2015/internal/db"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	db.Init()
	defer db.DB.Close()
	r := mux.NewRouter()
	r.HandleFunc("/api/auth", httphandler.AuthHandler).Methods("POST")
	r.HandleFunc("/api/sendCoin", httphandler.SendCoinHandler).Methods("POST")
	r.HandleFunc("/api/info", httphandler.InfoHandler).Methods("GET")
	r.HandleFunc("/api/buy/{item}", httphandler.BuyItemHandler).Methods("POST")

	// TODO: пробросить адрес сервера через конфиги
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
