package storage

import "time"

type Event struct {
	ID          string
	Title       string
	Start       time.Time
	Finish      time.Time
	Description string
	UserID      string
}

type User struct {
	ID   string
	Name string
}
