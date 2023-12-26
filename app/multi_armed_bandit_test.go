package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanSkripnikov/golang_otus_project/components"
)

func TestGetRating(t *testing.T) {
	expected := 1.924720344358278

	averageRating := 1.34
	allShowsBanner := float64(23)
	allShows := float64(51)

	result := components.GetRating(averageRating, allShowsBanner, allShows)
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func TestGetAllEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) from events WHERE type =").WithArgs("show")

	components.SetDatabase(db)
	_ = components.GetAllEvents("show")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBannerEvents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) as cnt from events WHERE banner_id = (.+) AND group_id = (.+) AND slot_id = (.+) AND type = (.+)").WithArgs(1, 1, 1, "show")

	components.SetDatabase(db)
	_ = components.GetBannerEvents(1, 1, 1, "show")

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBannersForSlot(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT banner_id from relations_banner_slot WHERE slot_id =").WithArgs(1)

	components.SetDatabase(db)
	_, err = components.GetBannersForSlot(1)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetBannerRatings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	banners := []int{1, 2}

	mock.ExpectQuery("SELECT (.+) as cnt from events WHERE banner_id = (.+) AND group_id = (.+) AND slot_id = (.+) AND type = (.+)").WithArgs(1, 1, 1, "show")
	mock.ExpectQuery("SELECT (.+) as cnt from events WHERE banner_id = (.+) AND group_id = (.+) AND slot_id = (.+) AND type = (.+)").WithArgs(1, 1, 1, "click")
	mock.ExpectQuery("SELECT (.+) from events WHERE type =").WithArgs("show")

	mock.ExpectQuery("SELECT (.+) as cnt from events WHERE banner_id = (.+) AND group_id = (.+) AND slot_id = (.+) AND type = (.+)").WithArgs(2, 1, 1, "show")
	mock.ExpectQuery("SELECT (.+) as cnt from events WHERE banner_id = (.+) AND group_id = (.+) AND slot_id = (.+) AND type = (.+)").WithArgs(2, 1, 1, "click")
	mock.ExpectQuery("SELECT (.+) from events WHERE type =").WithArgs("show")

	components.SetDatabase(db)
	_ = components.GetBannerRatings(banners, 1, 1)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
