package components

import (
	"database/sql"
	"fmt"
	"math"
	"sort"

	"github.com/IvanSkripnikov/golang_otus_project/database"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
	"github.com/IvanSkripnikov/golang_otus_project/models"
)

var db *sql.DB

func init() {
	SetDatabase(database.DB)
}

func SetDatabase(database *sql.DB) {
	db = database
}

func GetNeedBanner(slotID, groupID int) int {
	var resultBannerID int

	// находим баннеры для данного слота
	bannersForSlot, err := GetBannersForSlot(slotID)
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
		allShowsBanner := float64(GetBannerEvents(bannerID, groupID, slotID, "show"))
		allClickBanner := float64(GetBannerEvents(bannerID, groupID, slotID, "click"))
		allShows := float64(GetAllEvents("show"))

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

func GetAllEvents(eventType string) int {
	query := "SELECT COUNT(*) from events WHERE type = ?"
	rows, err := db.Query(query, eventType)
	if err != nil {
		return 0
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	count := 0

	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0
		}
	}

	return count
}

func GetBannerEvents(bannerID, groupID, slotID int, eventType string) int {
	query := "SELECT COUNT(*) as cnt from events WHERE banner_id = ? AND group_id = ? AND slot_id = ? AND type = ?"
	rows, err := db.Query(query, bannerID, groupID, slotID, eventType)
	if err != nil {
		return 0
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	count := 0

	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0
		}
	}

	return count
}

func GetBannersForSlot(slotID int) ([]int, error) {
	query := "SELECT banner_id from relations_banner_slot WHERE slot_id = ?"
	rows, err := db.Query(query, slotID)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	banners := make([]int, 0)
	banner := 0
	for rows.Next() {
		if err = rows.Scan(&banner); err != nil {
			logger.SendToErrorLog(err.Error())
			continue
		}
		banners = append(banners, banner)
	}

	return banners, nil
}
