package models

import "time"

type Test struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
