package database

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type RatingReviewRepository interface {
	GetAllOfProduct(
		c *gin.Context,
		db Querier,
		productID int32,
	) ([]RatingsAndReviewDetails, error)
}

type ratingReviewRepo struct{}

func NewRatingReviewRepository() RatingReviewRepository {
	return &ratingReviewRepo{}
}

type RatingsAndReviewDetails struct {
	ID        int32       `json:"id"`
	Rating    int32       `json:"rating"`
	Review    string      `json:"review"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	UserImage pgtype.Text `json:"userImage"`
	UserID    int32       `json:"userId"`
}

func (repo *ratingReviewRepo) GetAllOfProduct(
	c *gin.Context,
	db Querier,
	productID int32,
) ([]RatingsAndReviewDetails, error) {
	query := `
		SELECT 
			rr.id,
			rr.rating,
			rr.review,
			rr.created_at,
			rr.updated_at,
			u.first_name,
			u.last_name,
			u.image,
			u.id
		FROM rating_review rr
		JOIN users u ON rr.user_id = u.id
		WHERE rr.product_id = $1
		ORDER BY rr.created_at DESC
	`

	rows, err := db.Query(c, query, productID)
	if err != nil {
		return nil, Parse(err, "Rating Review", "GetAllOfProduct")
	}
	defer rows.Close()

	var rrs []RatingsAndReviewDetails
	for rows.Next() {
		var rr RatingsAndReviewDetails
		err = rows.Scan(
			&rr.ID,
			&rr.Rating,
			&rr.Review,
			&rr.CreatedAt,
			&rr.UpdatedAt,
			&rr.FirstName,
			&rr.LastName,
			&rr.UserImage,
			&rr.UserID,
		)
		if err != nil {
			return nil, Parse(err, "Rating Review", "GetAllOfProduct")
		}
		rrs = append(rrs, rr)
	}

	if err = rows.Err(); err != nil {
		return nil, Parse(err, "Rating Review", "GetAllOfProduct")
	}

	return rrs, nil
}
