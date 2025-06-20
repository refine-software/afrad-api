package database

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/models"
)

type RatingReviewRepository interface {
	// This method will create a Review, with the following data:
	// rating, review, user_id, product_id.
	Create(
		c *gin.Context,
		db Querier,
		rr *models.RatingReview,
	) error

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

func (repo *ratingReviewRepo) Create(
	c *gin.Context,
	db Querier,
	rr *models.RatingReview,
) error {
	query := `
		INSERT INTO rating_review (rating, review, user_id, product_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := db.QueryRow(c, query, rr.Rating, rr.Review, rr.UserID, rr.ProductID).
		Scan(&rr.ID, &rr.CreatedAt, &rr.UpdatedAt)
	if err != nil {
		return Parse(err, "Rating Review", "Create")
	}

	return nil
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
