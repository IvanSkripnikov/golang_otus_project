package controllers

import (
	"fmt"
	"net/http"

	"github.com/IvanSkripnikov/golang_otus_project/logger"
)

func HelloPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		_, err := fmt.Fprint(w, "{\"message\": \"Hello dear friend! Welcome!\"}")
		if err != nil {
			logger.SendToErrorLog(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
