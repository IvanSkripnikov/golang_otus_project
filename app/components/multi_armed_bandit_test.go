package components

import (
	"testing"

	"github.com/IvanSkripnikov/golang_otus_project/models"
)

func TestGetRating(t *testing.T) {
	expected := 1.924720344358278

	averageRating := 1.34

	allShowsBanner := float64(23)

	allShows := float64(51)

	result := GetRating(averageRating, allShowsBanner, allShows)

	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestGetBannerRatings(t *testing.T) {
	var allShows float64 = 128

	bannersStatistics := []models.BannerStats{
		{BannerID: 1, AllClickBanner: 1, AllShowsBanner: 10},

		{BannerID: 3, AllClickBanner: 5, AllShowsBanner: 20},
	}

	result := GetBannerRatings(allShows, bannersStatistics)

	if len(result) != 2 {
		t.Error("error get banners rating")
	}

	success := 0

	for _, v := range result {
		if (v.BannerID == 1 || v.BannerID == 3) && v.Value > 0 {
			success++
		}
	}

	if success != 2 {
		t.Error("unexpected result")
	}
}

func TestGetNeedBanner(t *testing.T) {
	var allShows float64 = 200

	bannersStatistics := []models.BannerStats{
		{BannerID: 1, AllClickBanner: 1, AllShowsBanner: 10},

		{BannerID: 2, AllClickBanner: 2, AllShowsBanner: 20},

		{BannerID: 3, AllClickBanner: 5, AllShowsBanner: 30},
	}

	bannerID := GetNeedBanner(allShows, bannersStatistics)

	if bannerID != 1 {
		t.Errorf("Expected %v but got %v", 1, bannerID)
	}
}
