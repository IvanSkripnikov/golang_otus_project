package components

import (
	"fmt"
	"math"
	"sort"

	"github.com/IvanSkripnikov/golang_otus_project/database"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
	"github.com/IvanSkripnikov/golang_otus_project/models"
)

func GetNeedBanner(slotID, groupID int) int {
	var resultBannerID int

	// находим баннеры для данного слота
	bannersForSlot, err := database.GetBannersForSlot(slotID)
	if err != nil {
		logger.SendToFatalLog("error while search banners.")
	}

	// получаем рейтинги по баннерам
	rateBanners := GetBannerRatings(bannersForSlot, groupID, slotID)

	resultBannerID = rateBanners[0].BannerID

	return resultBannerID
}

func GetRating(averageRating float64, currentCount float64, allCounts float64) float64 {
	return averageRating + math.Sqrt((2*math.Log(allCounts))/currentCount)
}

func GetBannerRatings(bannersForSlot []int, groupID, slotID int) []models.Rating {
	rateBanners := make([]models.Rating, len(bannersForSlot))

	var averageRating, rate float64
	for k, bannerID := range bannersForSlot {
		allShowsBanner := float64(database.GetBannerEvents(bannerID, groupID, slotID, "show"))
		allClickBanner := float64(database.GetBannerEvents(bannerID, groupID, slotID, "click"))
		allShows := float64(database.GetAllEvents("show"))

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
			BannerID: bannerID,
			Value:    rate,
		}

		rateBanners[k] = rating
	}

	// сортируем итоговый набор рейтингов
	sort.Slice(rateBanners, func(i, j int) bool {
		return rateBanners[i].Value > rateBanners[j].Value
	})

	fmt.Println(rateBanners)

	return rateBanners
}