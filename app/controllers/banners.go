package controllers

import (
	"database/sql"
	"net/http"

	"github.com/IvanSkripnikov/golang_otus_project/database"
	"github.com/IvanSkripnikov/golang_otus_project/helpers"
)

var db *sql.DB

func init() {
	SetDatabase(database.DB)
}

func SetDatabase(database *sql.DB) {
	db = database
}

func BannersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		helpers.GetAllBanners(w, r, db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func BannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		helpers.GetBanner(w, r, db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AddBannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		helpers.AddBannerToSlot(w, r, db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RemoveBannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		helpers.RemoveBannerFromSlot(w, r, db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetBannerForShowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		helpers.GetBannerForShow(w, r, db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		helpers.EventClick(w, r, db)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
