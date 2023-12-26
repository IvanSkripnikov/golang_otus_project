package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanSkripnikov/golang_otus_project/controllers"
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

func TestBanners(t *testing.T) {
	handler := GetHTTPHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	controllers.SetDatabase(db)

	rows := sqlmock.NewRows([]string{"1", "moscow aged", "Московское долголение", "08.06.2018 17:30", "1"})
	mock.ExpectQuery("SELECT (.+) from banners").WillReturnRows(rows)

	e := httpexpect.Default(t, server.URL)

	e.GET("/banners").
		Expect().
		Status(http.StatusOK).JSON()

	e.POST("/banners").Expect().Status(http.StatusMethodNotAllowed)
	e.PUT("/banners").Expect().Status(http.StatusMethodNotAllowed)
	e.PATCH("/banners").Expect().Status(http.StatusMethodNotAllowed)
	e.DELETE("/banners").Expect().Status(http.StatusMethodNotAllowed)
	e.HEAD("/banners").Expect().Status(http.StatusMethodNotAllowed)
}
