package models

type Banners struct {
	Bannerx []Banner `json:"bannerx"`
}
type Banner struct {
	Url string
}
