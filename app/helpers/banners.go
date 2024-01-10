package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/IvanSkripnikov/golang_otus_project/components"
	"github.com/IvanSkripnikov/golang_otus_project/database"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
	"github.com/IvanSkripnikov/golang_otus_project/models"
	"github.com/IvanSkripnikov/golang_otus_project/queue"
)

func GetAllBanners(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "SELECT * from banners"
	rows, err := database.DB.Query(query)

	if checkError(w, err) {
		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var banners []models.Banner

	for rows.Next() {
		banner := models.Banner{}

		if err = rows.Scan(&banner.ID, &banner.Title, &banner.Body, &banner.CreatedAt, &banner.Active); err != nil {
			logger.SendToErrorLog(err.Error())

			continue
		}
		banners = append(banners, banner)
	}

	var buf bytes.Buffer

	je := json.NewEncoder(&buf)
	err = je.Encode(&banners)

	if checkError(w, err) {
		return
	}

	writeSuccess(w, buf.String())
}

func GetBanner(w http.ResponseWriter, r *http.Request) {
	var banner models.Banner

	banner.ID, _ = getIDFromRequestString(r.URL.Path)

	if banner.ID == 0 {
		wrongParamsResponse(w)
		return
	}

	query := "SELECT * from banners WHERE id = ?"
	rows, err := database.DB.Prepare(query)

	if checkError(w, err) {
		return
	}

	defer rows.Close()

	err = rows.QueryRow(banner.ID).Scan(&banner.ID, &banner.Title, &banner.Body, &banner.CreatedAt, &banner.Active)

	if err != nil {
		logger.SendToErrorLog(err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "{ \"message\": \"Not Found\"}")

		return
	}

	var buf bytes.Buffer
	je := json.NewEncoder(&buf)
	err = je.Encode(&banner)

	if checkError(w, err) {
		return
	}

	writeSuccess(w, buf.String())
}

func AddBannerToSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)
	slotID, okSlot := params["slot"]
	bannerID, okBanner := params["banner"]

	if !okSlot || !okBanner || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	query := "INSERT INTO relations_banner_slot (banner_id, slot_id) VALUES (?, ?)"
	rows, err := database.DB.Query(query, bannerID, slotID)

	if checkError(w, err) {
		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	writeSuccess(w, "{\"message\": \"Successfully added!\"}")
}

func RemoveBannerFromSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)
	slotID, okSlot := params["slot"]
	bannerID, okBanner := params["banner"]

	if !okSlot || !okBanner || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	query := "DELETE FROM relations_banner_slot WHERE banner_id=? AND slot_id=?"
	rows, err := database.DB.Query(query, bannerID, slotID)

	if checkError(w, err) {
		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	writeSuccess(w, "{\"message\": \"Successfully removed!\"}")
}

func GetBannerForShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)
	slotID, okSlot := params["slot"]
	groupID, okGroup := params["group"]

	if !okSlot || !okGroup || resultString != "" {
		wrongParamsResponse(w)

		return
	}

	allShows := float64(GetAllEvents("show"))
	bannersStatistics := getBannersStatistics(slotID, groupID)

	// получаем id баннера
	bannerID := components.GetNeedBanner(allShows, bannersStatistics)

	// записываем событие просмотра
	query := "INSERT INTO events (type, banner_id, slot_id, group_id) VALUES (?, ?, ?, ?)"
	rows, err := database.DB.Query(query, "show", bannerID, slotID, groupID)

	if checkError(w, err) {
		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	// отправляем событие в rabbitMQ

	go func() {
		queue.SendEventToQueue("show", bannerID, slotID, groupID)
	}()

	_, err = fmt.Fprint(w, strconv.Itoa(bannerID))

	if err != nil {
		logger.SendToErrorLog(err.Error())

		return
	}
}

func EventClick(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params, resultString := getParamsFromQueryString(r.URL.Path)
	slotID, okSlot := params["slot"]
	groupID, okGroup := params["group"]
	bannerID, okBanner := params["banner"]

	if !okSlot || !okGroup || !okBanner || resultString != "" {
		wrongParamsResponse(w)
		return
	}

	query := "INSERT INTO events (type, banner_id, slot_id, group_id) VALUES (?, ?, ?, ?)"
	rows, err := database.DB.Query(query, "click", bannerID, slotID, groupID)

	// отправляем событие в rabbitMQ
	go func() {
		queue.SendEventToQueue("click", bannerID, slotID, groupID)
	}()

	if checkError(w, err) {
		return
	}

	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	fmt.Fprint(w, r)
}

func GetAllEvents(eventType string) int {
	query := "SELECT COUNT(*) from events WHERE type = ?"

	rows, err := database.DB.Query(query, eventType)
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

	rows, err := database.DB.Query(query, slotID)
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

func GetBannerEvents(bannerID, groupID, slotID int, eventType string) int {
	query := "SELECT COUNT(*) as cnt from events WHERE banner_id = ? AND group_id = ? AND slot_id = ? AND type = ?"

	rows, err := database.DB.Query(query, bannerID, groupID, slotID, eventType)
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

// -------------PRIVATE----------------------

func getBannersStatistics(slotID, groupID int) []models.BannerStats {
	// находим баннеры для данного слота
	bannersForSlot, err := GetBannersForSlot(slotID)
	if err != nil {
		logger.SendToFatalLog("error while search banners.")
	}

	bannerStats := make([]models.BannerStats, len(bannersForSlot))

	for k, bannerID := range bannersForSlot {
		bannerStats[k] = models.BannerStats{
			BannerID:       bannerID,
			AllShowsBanner: float64(GetBannerEvents(bannerID, groupID, slotID, "show")),
			AllClickBanner: float64(GetBannerEvents(bannerID, groupID, slotID, "click")),
		}
	}

	return bannerStats
}

func checkError(w http.ResponseWriter, err error) bool {
	if err != nil {
		logger.SendToErrorLog(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return true
	}

	return false
}

func writeSuccess(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)

	_, err := fmt.Fprint(w, message)
	if err != nil {
		logger.SendToErrorLog(fmt.Sprintf("write success error %s", err.Error()))

		return
	}
}

func wrongParamsResponse(w http.ResponseWriter) {
	resultString := "{\"message\": \"Invalid request GetHandler\"}"
	fmt.Fprint(w, resultString)
	w.WriteHeader(http.StatusBadRequest)
}

func getIDFromRequestString(url string) (int, error) {
	vars := strings.Split(url, "/")

	return strconv.Atoi(vars[len(vars)-1])
}

func getParamsFromQueryString(url string) (map[string]int, string) {
	resultMap := map[string]int{}
	outMessage := ""
	queryParams := strings.Split(url, "/")
	params := strings.Split(queryParams[len(queryParams)-1], "&")

	if len(params) == 1 {
		outMessage = "not all params is set"
		return resultMap, outMessage
	}

	for _, v := range params {
		pair := strings.Split(v, "=")

		if len(pair) == 1 {
			outMessage = "incorrect params value: " + v
			return resultMap, outMessage
		}

		resultMap[pair[0]], _ = strconv.Atoi(pair[1])
	}

	return resultMap, outMessage
}
