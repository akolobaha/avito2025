package httphandler

import "net/http"

func BuyItemHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"message": "Authentication successful"}`))
	if err != nil {
		return
	}
}
