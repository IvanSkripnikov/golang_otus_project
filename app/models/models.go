package models

type Banner struct {
	ID        int
	Title     string
	Body      string
	CreatedAt string
	Active    bool
}

type Rating struct {
	BannerID int
	Value    float64
}
