package components

import (
	"math"
	"sort"

	"github.com/IvanSkripnikov/golang_otus_project/models"
)

func GetNeedBanner(allShows float64, bannersStatistics []models.BannerStats) int {
	var resultBannerID int

	// получаем рейтинги по баннерам

	rateBanners := GetBannerRatings(allShows, bannersStatistics)

	resultBannerID = rateBanners[0].BannerID

	return resultBannerID
}

func GetRating(averageRating float64, currentCount float64, allCounts float64) float64 {
	return averageRating + math.Sqrt((2*math.Log(allCounts))/currentCount)
}

func GetBannerRatings(allShows float64, bannersStatistics []models.BannerStats) []models.Rating {
	rateBanners := make([]models.Rating, len(bannersStatistics))

	var averageRating, rate float64

	for k, bannerStat := range bannersStatistics {
		allShowsBanner := bannerStat.AllShowsBanner

		allClickBanner := bannerStat.AllClickBanner

		// находим средний рейтинг баннера

		if allClickBanner == 0 || allShowsBanner == 0 {
			averageRating = 0
		} else {
			averageRating = allClickBanner / allShowsBanner
		}

		// считаем рейтинг баннера

		if allShowsBanner == 0 {
			rate = 0
		} else {
			rate = GetRating(averageRating, allShowsBanner, allShows)
		}

		rating := models.Rating{
			BannerID: bannerStat.BannerID,

			Value: rate,
		}

		rateBanners[k] = rating
	}

	// сортируем итоговый набор рейтингов

	sort.Slice(rateBanners, func(i, j int) bool {
		return rateBanners[i].Value > rateBanners[j].Value
	})

	return rateBanners
}
