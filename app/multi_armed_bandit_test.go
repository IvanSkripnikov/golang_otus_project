package main

import (
	"testing"

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
