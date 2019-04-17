package main

import "time"

// Question fd
type Question struct {
	ID        int64 `gorm:"primary_key"`
	EventID   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Content   string
	Answers   string
}
