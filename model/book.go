package model

type Book struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Year   uint   `json:"year" validate:"required"`
}