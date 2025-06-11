package database

type RatingReviewRepository interface{}

type ratingReviewRepo struct{}

func NewRatingReviewRepository() RatingReviewRepository {
	return &ratingReviewRepo{}
}
