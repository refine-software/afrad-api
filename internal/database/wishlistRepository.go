package database

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/models"
)

type WishlistRepository interface {
	GetAllOfUser(c *gin.Context, db Querier, userID int32) ([]wishlist, error)

	// This method creates a wishlist record,
	// columns required: product_id, user_id
	Create(*gin.Context, Querier, *models.Wishlist) error

	Delete(c *gin.Context, db Querier, wishlistID int32) error
}

type wishlistRepo struct{}

func NewWishlistRepository() WishlistRepository {
	return &wishlistRepo{}
}

type wishlist struct {
	ID               int32     `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UserID           int32     `json:"userId"`
	ProductID        int32     `json:"productId"`
	ProductName      string    `json:"productName"`
	ProductThumbnail string    `json:"productThumbnail"`
}

func (repo *wishlistRepo) GetAllOfUser(
	c *gin.Context,
	db Querier,
	userID int32,
) ([]wishlist, error) {
	query := `
		SELECT w.id, w.created_at, w.user_id, w.product_id, p.name, p.thumbnail
		FROM wishlists w
		JOIN products p ON w.product_id = p.id
		WHERE user_id = $1
	`

	rows, err := db.Query(c, query, userID)
	if err != nil {
		return nil, Parse(err, "Wishlist", "GetAllOfUser", make(Constraints))
	}
	defer rows.Close()

	var ws []wishlist
	for rows.Next() {
		var w wishlist
		err = rows.Scan(
			&w.ID,
			&w.CreatedAt,
			&w.UserID,
			&w.ProductID,
			&w.ProductName,
			&w.ProductThumbnail,
		)
		if err != nil {
			return nil, Parse(err, "Wishlist", "GetAllOfUser", make(Constraints))
		}

		ws = append(ws, w)
	}

	if rows.Err() != nil {
		return nil, Parse(err, "Wishlist", "GetAllOfUser", make(Constraints))
	}

	return ws, nil
}

func (repo *wishlistRepo) Create(c *gin.Context, db Querier, r *models.Wishlist) error {
	query := `
		INSERT INTO wishlists(user_id, product_id)
		VALUES($1, $2)
	`

	_, err := db.Exec(c, query, r.UserID, r.ProductID)
	if err != nil {
		return Parse(err, "Wishlist", "Create", Constraints{
			NotNullViolationCode:    "product",
			ForeignKeyViolationCode: "product",
			UniqueViolationCode:     "product",
		})
	}

	return nil
}

func (repo *wishlistRepo) Delete(c *gin.Context, db Querier, wishlistID int32) error {
	query := `
		DELETE FROM wishlists
		WHERE id = $1
	`

	result, err := db.Exec(c, query, wishlistID)
	if err != nil {
		return Parse(err, "Wishlist", "Delete", make(Constraints))
	}

	if result.RowsAffected() == 0 {
		return Parse(pgx.ErrNoRows, "Wishlist", "Delete", make(Constraints))
	}

	return nil
}
