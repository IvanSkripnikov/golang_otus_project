package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IvanSkripnikov/golang_otus_project/controllers"
	"github.com/IvanSkripnikov/golang_otus_project/database"
	"github.com/IvanSkripnikov/golang_otus_project/logger"
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
	handler := GetHttpHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners/1").
		Expect().
		Status(http.StatusOK).JSON().IsObject()
}

func TestBanners(t *testing.T) {
	handler := GetHttpHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners").
		Expect().
		Status(http.StatusOK).JSON().Array().NotEmpty()
}

func TestGetBannerForShow(t *testing.T) {
	countBefore := database.GetAllEvents("show")
	handler := GetHttpHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.GET("/get_banner_for_show/slot=1&group=1").
		Expect().
		Status(http.StatusOK).Raw()

	countAfter := database.GetAllEvents("show")

	if countBefore+1 != countAfter {
		t.Errorf("Not increment banner shows!")
	}
}

func TestEventClick(t *testing.T) {
	countBefore := database.GetAllEvents("click")

	handler := GetHttpHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/event_click/slot=1&group=1&banner=1").
		Expect().
		Status(http.StatusOK).Raw()

	countAfter := database.GetAllEvents("click")

	if countBefore+1 != countAfter {
		t.Errorf("Not increment banner shows!")
	}
}

func TestRemoveBanner(t *testing.T) {
	handler := GetHttpHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	tx, err := database.DB.BeginTx(context.Background(), nil)
	if err != nil {
		logger.SendToFatalLog(err.Error())
	}
	defer tx.Rollback()
	// далее - обычная работа как с *sql.DB

	e.POST("/remove_banner_from_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).Raw()
}

func TestAddBanner(t *testing.T) {
	handler := GetHttpHandler()
	server := httptest.NewServer(handler)
	defer server.Close()

	e := httpexpect.Default(t, server.URL)

	e.POST("/add_banner_to_slot/slot=1&banner=1").
		Expect().
		Status(http.StatusOK).Raw()
}
