package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IvanSkripnikov/golang_otus_project/controllers"
	"github.com/IvanSkripnikov/golang_otus_project/helpers"
	"github.com/gavv/httpexpect/v2"
)

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

func TestBanner(t *testing.T) {
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

func TestGetBannerForShow(t *testing.T) {
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

func TestEventClick(t *testing.T) {
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

func TestRemoveBanner(t *testing.T) {
	handler := GetHTTPHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/remove_banner_from_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON()

	e.GET("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.HEAD("/remove_banner_from_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
}

func TestAddBanner(t *testing.T) {
	handler := GetHTTPHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/add_banner_to_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).JSON()

	e.GET("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
	e.HEAD("/add_banner_to_slot/slot=1&banner=1").Expect().Status(http.StatusMethodNotAllowed)
}
