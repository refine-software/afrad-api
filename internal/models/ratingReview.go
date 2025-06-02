package models

import "time"

type RatingReview struct {
	ID        int
	Rating    int
	Review    string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    int
	ProductID int
}
