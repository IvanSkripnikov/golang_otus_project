package main

import (
	"testing"

	"github.com/IvanSkripnikov/golang_otus_project/components"
)

func TestGetRating(t *testing.T) {
	expected := 1.924720344358278

	averageRating := 1.34

	allShowsBanner := float64(23)

	allShows := float64(51)

	result := components.GetRating(averageRating, allShowsBanner, allShows)

	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestGetAllEvents(t *testing.T) {
	countEvents := components.GetAllEvents("show")
	if countEvents == 0 {
		t.Error("error check events count")
	}
}

func TestGetBannerEvents(t *testing.T) {
	countBanners := components.GetBannerEvents(1, 1, 1, "show")
	if countBanners == 0 {
		t.Error("error check banner events count")
	}
}

func TestGetBannersForSlot(t *testing.T) {
	banners, err := components.GetBannersForSlot(1)

	if err != nil {
		t.Errorf("error while get banners for slot: %v", err)
	}
	if len(banners) == 0 {
		t.Error("error check banners count")
	}
}

func TestGetBannerRatings(t *testing.T) {
	banners := []int{1, 3}
	result := components.GetBannerRatings(banners, 1, 1)

	if len(result) != 2 {
		t.Error("error get banners rating")
	}

	if result[0].BannerID != 1 || result[0].Value == 0 || result[1].BannerID != 3 || result[1].Value == 0 {
		t.Error("unexpected result")
	}
}
