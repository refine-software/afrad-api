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

	Get(c *gin.Context, db Querier, reviewID int32) (*models.RatingReview, error)

	GetAllOfProduct(
		c *gin.Context,
		db Querier,
		productID int32,
	) ([]RatingsAndReviewDetails, error)

	GetAllOfUser(c *gin.Context, db Querier, userID int32) ([]models.RatingReview, error)

	// This method updates the RatingReview model,
	// Required columns: rating, review.
	// By: id.
	Update(c *gin.Context, db Querier, rr *models.RatingReview) error

	// This method will delete a review from the database
	// by the review id.
	Delete(c *gin.Context, db Querier, reviewID int32) error
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
		return Parse(err, "Rating Review", "Create", map[string]string{
			UniqueViolationCode:     "review",
			ForeignKeyViolationCode: "product",
		})
	}

	return nil
}

func (repo *ratingReviewRepo) Get(
	c *gin.Context,
	db Querier,
	reviewID int32,
) (*models.RatingReview, error) {
	query := `
		SELECT id, rating, review, created_at, updated_at, user_id, product_id
		FROM rating_review
		WHERE id = $1
	`

	var rr models.RatingReview
	err := db.QueryRow(c, query, reviewID).
		Scan(&rr.ID, &rr.Rating, &rr.Review, &rr.CreatedAt, &rr.UpdatedAt, &rr.UserID, &rr.ProductID)
	if err != nil {
		return nil, Parse(err, "Rating Review", "Get", make(Constraints))
	}

	return &rr, nil
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
		return nil, Parse(err, "Rating Review", "GetAllOfProduct", make(Constraints))
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
			return nil, Parse(err, "Rating Review", "GetAllOfProduct", make(Constraints))
		}
		rrs = append(rrs, rr)
	}

	if err = rows.Err(); err != nil {
		return nil, Parse(err, "Rating Review", "GetAllOfProduct", make(Constraints))
	}

	return rrs, nil
}

func (repo *ratingReviewRepo) GetAllOfUser(
	c *gin.Context,
	db Querier,
	userID int32,
) ([]models.RatingReview, error) {
	query := `
		SELECT id, rating, review, created_at, updated_at, user_id, product_id
		FROM rating_review
		WHERE user_id = $1
	`

	var rrs []models.RatingReview

	rows, err := db.Query(c, query, userID)
	if err != nil {
		return nil, Parse(err, "Rating Review", "GetAllOfUser", make(Constraints))
	}
	defer rows.Close()

	for rows.Next() {
		var rr models.RatingReview
		err = rows.Scan(
			&rr.ID,
			&rr.Rating,
			&rr.Review,
			&rr.CreatedAt,
			&rr.UpdatedAt,
			&rr.UserID,
			&rr.ProductID,
		)
		if err != nil {
			return nil, Parse(err, "Rating Review", "GetAllOfUser", make(Constraints))
		}

		rrs = append(rrs, rr)
	}

	return rrs, nil
}

func (repo *ratingReviewRepo) Update(c *gin.Context, db Querier, rr *models.RatingReview) error {
	query := `
		UPDATE rating_review
		SET rating = $2, review = $3
		WHERE id = $1
		RETURNING updated_at
	`

	err := db.QueryRow(c, query, rr.ID, rr.Rating, rr.Review).Scan(&rr.UpdatedAt)
	if err != nil {
		return Parse(err, "Rating Review", "Update", Constraints{
			NotNullViolationCode:          "rating",
			CheckViolationCode:            "rating",
			StringDataRightTruncationCode: "review",
		})
	}

	return nil
}

func (repo *ratingReviewRepo) Delete(c *gin.Context, db Querier, reviewID int32) error {
	query := `
		DELETE FROM rating_review
		WHERE id = $1
	`

	_, err := db.Exec(c, query, reviewID)
	if err != nil {
		return Parse(err, "Rating Review", "Delete", make(Constraints))
	}

	return nil
}
