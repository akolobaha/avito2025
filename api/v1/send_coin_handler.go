package httphandler

import "net/http"

func SendCoinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"message": "Send coin"}`))
	if err != nil {
		return
	}
}
