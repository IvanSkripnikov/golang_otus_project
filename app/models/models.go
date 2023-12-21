package models

type Banner struct {
	Id        int
	Title     string
	Body      string
	CreatedAt string
	Active    bool
}

type Rating struct {
	BannerId int
	Value    float64
}
