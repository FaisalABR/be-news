package entity

import "time"

type ContentEntity struct {
	ID          int64
	Title       string
	Excerpt     string
	Description string
	Image       string
	Status      string
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CategoryID  int64
	CreatedByID int64
	Category    CategoryEntity
	User        UserEntity
}

type QueryString struct {
	Page       int
	Limit      int
	OrderBy    string
	OrderType  string
	Search     string
	Status     string
	CategoryID int
}
