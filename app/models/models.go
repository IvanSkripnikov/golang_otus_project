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

type Message struct {
	Type     string
	BannerID int
	SlotID   int
	GroupID  int
}

type BannerStats struct {
	BannerID       int
	AllShowsBanner float64
	AllClickBanner float64
}
