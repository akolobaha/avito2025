package httphandler

import "net/http"

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"message": "Info handler"}`))
	if err != nil {
		return
	}
}
