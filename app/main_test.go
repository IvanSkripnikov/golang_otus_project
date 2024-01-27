package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IvanSkripnikov/golang_otus_project/controllers"
	"github.com/IvanSkripnikov/golang_otus_project/database"
	"github.com/IvanSkripnikov/golang_otus_project/helpers"
	"github.com/gavv/httpexpect/v2"
)

func init() {
	database.InitDataBase("localhost")
}

func TestRoot(t *testing.T) {
	expected := "{\"message\": \"Hello dear friend! Welcome!\"}"

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	w := httptest.NewRecorder()

	controllers.HelloPageHandler(w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected root message but got %v", string(data))
	}
}

func TestBannerSuccess(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners/1").
		Expect().
		Status(http.StatusOK).JSON().IsObject()

	e.POST("/banners/1").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/banners/1").Expect().Status(http.StatusMethodNotAllowed)

	e.PATCH("/banners/1").Expect().Status(http.StatusMethodNotAllowed)

	e.DELETE("/banners/1").Expect().Status(http.StatusMethodNotAllowed)

	e.HEAD("/banners/1").Expect().Status(http.StatusMethodNotAllowed)
}

func TestBannerNotFound(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners/-1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestBanners(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners").
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty()

	e.POST("/banners").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/banners").Expect().Status(http.StatusMethodNotAllowed)

	e.PATCH("/banners").Expect().Status(http.StatusMethodNotAllowed)

	e.DELETE("/banners").Expect().Status(http.StatusMethodNotAllowed)

	e.HEAD("/banners").Expect().Status(http.StatusMethodNotAllowed)
}

func TestGetBannerForShowSuccess(t *testing.T) {
	countBefore := helpers.GetAllEvents("show")

	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/get_banner_for_show/slot=1&group=1").
		Expect().
		Status(http.StatusOK).JSON()

	countAfter := helpers.GetAllEvents("show")

	if countBefore+1 != countAfter {
		t.Errorf("Not increment banner shows!")
	}

	e.POST("/get_banner_for_show/slot=1&group=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/get_banner_for_show/slot=1&group=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PATCH("/get_banner_for_show/slot=1&group=1").Expect().Status(http.StatusMethodNotAllowed)

	e.DELETE("/get_banner_for_show/slot=1&group=1").Expect().Status(http.StatusMethodNotAllowed)

	e.HEAD("/get_banner_for_show/slot=1&group=1").Expect().Status(http.StatusMethodNotAllowed)
}

func TestGetBannerForShowFailureNoAssignedBannersForSlot(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/get_banner_for_show/slot=2&group=1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestGetBannerForShowFailureWrongSlot(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/get_banner_for_show/slot=-1&group=1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestGetBannerForShowFailureWrongGroup(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/get_banner_for_show/slot=1&group=-1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestEventClickSuccess(t *testing.T) {
	countBefore := helpers.GetAllEvents("click")

	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/event_click/slot=1&group=1&banner=1").
		Expect().
		Status(http.StatusOK)

	countAfter := helpers.GetAllEvents("click")

	if countBefore+1 != countAfter {
		t.Errorf("Not increment banner shows!")
	}

	e.GET("/event_click/slot=1&group=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/event_click/slot=1&group=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PATCH("/event_click/slot=1&group=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.DELETE("/event_click/slot=1&group=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.HEAD("/event_click/slot=1&group=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
}

func TestEventClickFailureWrongSlot(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/event_click/slot=-1&group=1&banner=1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestEventClickFailureWrongGroup(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/event_click/slot=1&group=-1&banner=1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestEventClickFailureWrongBanner(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/event_click/slot=1&group=1&banner=-1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestRemoveBannerSuccess(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	// сперва добавляем баннер, чтобы потом его удалить
	e.POST("/add_banner_to_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON()

	e.DELETE("/remove_banner_from_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON()

	e.GET("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PATCH("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.HEAD("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
}

func TestRemoveBannerFailureWrongBanner(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.DELETE("/remove_banner_from_slot/slot=1&banner=-1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestRemoveBannerFailureWrongSlot(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.DELETE("/remove_banner_from_slot/slot=-1&banner=1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestAddBannerSuccess(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/add_banner_to_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON()

	// убираем добавленный баннер
	e.DELETE("/remove_banner_from_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON()

	e.GET("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PUT("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.PATCH("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.DELETE("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)

	e.HEAD("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
}

func TestAddBannerFailureWrongBanner(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/add_banner_to_slot/slot=1&banner=-1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestAddBannerFailureWrongSlot(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/add_banner_to_slot/slot=-1&banner=1").
		Expect().
		Status(http.StatusNotFound).JSON().IsObject()
}

func TestAddBannerFailureExistsRelation(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)

	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	// добавляем первый баннер
	e.POST("/add_banner_to_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON().IsObject()

	// пробуем добавить такой же
	e.POST("/add_banner_to_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().IsObject()

	// убираем добавленный баннер
	e.DELETE("/remove_banner_from_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON().IsObject()
}
