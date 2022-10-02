package models

type User struct {
	Id           string `json:"id"`
	InstagramUID string `json:"igid"`
	Slug         string `json:"slug"`
}
