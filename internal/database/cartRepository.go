package database

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/models"
)

type CartRepository interface {
	// Create cart by user id
	Create(ctx *gin.Context, db Querier, userID int32) (int32, error)

	// Update cart by user id
	// required columns: total_price, quantity
	Update(ctx *gin.Context, db Querier, cart *models.Cart) error

	// Get cart id by user id
	GetIDByUserID(ctx *gin.Context, db Querier, userID int32) (int32, error)

	// Get cart by user id
	GetByUserID(ctx *gin.Context, db Querier, userID int32) (*models.Cart, error)

	// Delete cart by user id
	Delete(ctx *gin.Context, db Querier, userID int32) error
}

type cartRepo struct{}

func NewCartRepository() CartRepository {
	return &cartRepo{}
}

func (r *cartRepo) Create(ctx *gin.Context, db Querier, userID int32) (int32, error) {
	query := `
		INSERT INTO carts(total_price, quantity, user_id)
		VALUES(0, 0, $1)
		RETURNING id
	`
	var cartID int32
	err := db.QueryRow(ctx, query, userID).Scan(&cartID)
	if err != nil {
		return 0, Parse(
			err,
			"Cart",
			"Create",
			Constraints{UniqueViolationCode: "user_id", NotNullViolationCode: "user_id"},
		)
	}

	return cartID, nil
}

func (r *cartRepo) GetIDByUserID(ctx *gin.Context, db Querier, userID int32) (int32, error) {
	query := `
		SELECT id
		FROM carts
		WHERE user_id = $1
	`

	var cartID int32
	err := db.QueryRow(ctx, query, userID).Scan(&cartID)
	if err != nil {
		return 0, Parse(err, "Cart", "GetByID", make(Constraints))
	}
	return cartID, nil
}

func (r *cartRepo) Update(ctx *gin.Context, db Querier, cart *models.Cart) error {
	query := `
		UPDATE carts
		SET 
			total_price = $2,
			quantity = $3
		WHERE user_id = $1
	`

	_, err := db.Exec(ctx, query, cart.UserID, cart.TotalPrice, cart.Quantity)
	if err != nil {
		return Parse(err, "Cart", "Update", make(Constraints))
	}
	return nil
}

func (r *cartRepo) GetByUserID(
	ctx *gin.Context,
	db Querier,
	userID int32,
) (*models.Cart, error) {
	query := `
		SELECT id, total_price, quantity
		FROM carts
		WHERE user_id = $1
	`
	var cart models.Cart
	err := db.QueryRow(ctx, query, userID).Scan(&cart.ID, &cart.TotalPrice, &cart.Quantity)
	if err != nil {
		return nil, Parse(err, "Cart", "GetByUserID", make(Constraints))
	}
	return &cart, nil
}

func (r *cartRepo) Delete(ctx *gin.Context, db Querier, userID int32) error {
	query := `
		DELETE FROM carts
		WHERE user_id = $1
	`
	result, err := db.Exec(ctx, query, userID)
	if err != nil {
		return Parse(err, "Cart", "Delete", make(Constraints))
	}
	if result.RowsAffected() == 0 {
		return Parse(pgx.ErrNoRows, "Cart", "Delete", make(Constraints))
	}
	return nil
}
