package components

import (
	"app/database"
	"app/logger"
	"app/models"
	"fmt"
	"math"
	"sort"
)

func GetNeedBanner(slotId, groupId int) int {
	resultBannerId := 0

	// находим баннеры для данного слота
	bannersForSlot, err := database.GetBannersForSlot(slotId)
	if err != nil {
		logger.SendToFatalLog("error while search banners.")
	}

	// получаем рейтинги по баннерам
	rateBanners := GetBannerRatings(bannersForSlot, groupId, slotId)

	resultBannerId = rateBanners[0].BannerId

	return resultBannerId
}

func GetRating(averageRating float64, currentCount float64, allCounts float64) float64 {
	return averageRating + math.Sqrt((2*math.Log(allCounts))/currentCount)
}

func GetBannerRatings(bannersForSlot []int, groupId, slotId int) []models.Rating {
	rateBanners := make([]models.Rating, len(bannersForSlot))

	var averageRating, rate float64
	for k, bannerId := range bannersForSlot {
		allShowsBanner := float64(database.GetBannerEvents(bannerId, groupId, slotId, "show"))
		allClickBanner := float64(database.GetBannerEvents(bannerId, groupId, slotId, "click"))
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
			BannerId: bannerId,
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
