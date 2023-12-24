package controllers

import (
	"net/http"

	"github.com/IvanSkripnikov/golang_otus_project/helpers"
)

const (
	httpStatusGet  = "GET"
	httpStatusPOST = "POST"
)

func BannersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case httpStatusGet:
		helpers.GetAllBanners(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func BannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case httpStatusGet:
		helpers.GetBanner(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AddBannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case httpStatusPOST:
		helpers.AddBannerToSlot(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func RemoveBannerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case httpStatusPOST:
		helpers.RemoveBannerFromSlot(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func GetBannerForShowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case httpStatusGet:
		helpers.GetBannerForShow(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ClickHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case httpStatusPOST:
		helpers.EventClick(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
