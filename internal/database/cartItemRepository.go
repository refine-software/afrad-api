package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type CartItemRepository interface {
	// Create cart item.
	//
	// required columns: cart_id, product_id, total_price, quantity
	Create(ctx *gin.Context, db Querier, cartItem *models.CartItem) error

	// This method will get.
	//
	// The total_price of the whole cart as well as the whole quantity
	GetPriceQuantityByCartID(ctx *gin.Context, db Querier, cartID int32) (int, int, error)

	// Get cart Items by cart_id
	GetAll(ctx *gin.Context, db Querier, cartID int32) ([]GetCartItems, error)

	// Update cart item quantity by id
	Update(ctx *gin.Context, db Querier, id int32, quantity int) error

	// Check cart item Existence by id
	CheckExistence(ctx *gin.Context, db Querier, id int32) error

	// Delete cart item by id
	Delete(ctx *gin.Context, db Querier, id int32) error

	// Delete cart items by cart id
	DeleteByCartID(ctx *gin.Context, db Querier, cartID int32) error
}

type cartItemRepo struct{}

func NewCartItemRepository() CartItemRepository {
	return &cartItemRepo{}
}

func (r *cartItemRepo) Create(ctx *gin.Context, db Querier, cartItem *models.CartItem) error {
	query := `
		INSERT INTO cart_items(cart_id, product_id, total_price, quantity)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.Exec(
		ctx,
		query,
		cartItem.CartID,
		cartItem.ProductID,
		cartItem.TotalPrice,
		cartItem.Quantity,
	)
	if err != nil {
		return Parse(
			err,
			"Cart Item",
			"CreateCartItem",
			Constraints{
				UniqueViolationCode:     "product",
				NotNullViolationCode:    "all columns",
				ForeignKeyViolationCode: "cart_id, product_id",
			},
		)
	}
	return nil
}

func (r *cartItemRepo) GetPriceQuantityByCartID(
	ctx *gin.Context,
	db Querier,
	cartID int32,
) (int, int, error) {
	query := `
		SELECT
			COALESCE(SUM(product_variants.price * cart_items.quantity), 0) AS total_price,
			COALESCE(SUM(cart_items.quantity), 0) AS total_quantity
		FROM cart_items
		JOIN product_variants on cart_items.product_id = product_variants.id
		WHERE cart_id = $1
	`

	var totalPrice int
	var quantity int

	err := db.QueryRow(ctx, query, cartID).Scan(&totalPrice, &quantity)
	fmt.Println(err)
	if err != nil {
		return 0, 0, Parse(err, "CartItem", "GetPriceQuantityByCartID", make(Constraints))
	}

	return totalPrice, quantity, nil
}

type GetCartItems struct {
	ID           int32  `json:"id"`
	Quantity     int    `json:"quantity"`
	TotalPrice   int    `json:"totalPrice"`
	VariantID    int32  `json:"variantId"`
	ProductName  string `json:"productName"`
	ProductImg   string `json:"productImg"`
	ProductPrice int    `json:"productPrice"`
	ColorID      int32  `json:"colorId"`
	SizeID       int32  `json:"sizeId"`
}

func (r cartItemRepo) GetAll(
	ctx *gin.Context,
	db Querier,
	cartID int32,
) ([]GetCartItems, error) {
	query := `
		SELECT
			cart_items.id, 
			cart_items.quantity,
			cart_items.total_price,
			product_variants.id AS variant_id,
			products.name,
			products.thumbnail,
			product_variants.price,
			product_variants.color_id,
			product_variants.size_id
		FROM cart_items
		JOIN product_variants on product_variants.id = cart_items.product_id
		JOIN products on products.id = product_variants.product_id
		WHERE cart_items.cart_id = $1	
	`
	var cartItems []GetCartItems
	rows, err := db.Query(ctx, query, cartID)
	if err != nil {
		return nil, Parse(err, "Cart Item", "GetAll", make(Constraints))
	}
	defer rows.Close()

	for rows.Next() {
		var i GetCartItems
		err := rows.Scan(
			&i.ID,
			&i.Quantity,
			&i.TotalPrice,
			&i.VariantID,
			&i.ProductName,
			&i.ProductImg,
			&i.ProductPrice,
			&i.ColorID,
			&i.SizeID,
		)
		if err != nil {
			return nil, Parse(err, "Cart Item", "GetAll", make(Constraints))
		}
		cartItems = append(cartItems, i)
	}
	return cartItems, nil
}

func (r *cartItemRepo) Update(ctx *gin.Context, db Querier, id int32, quantity int) error {
	query := `
		UPDATE cart_items 
		SET quantity = $2 
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id, quantity)
	if err != nil {
		return Parse(err, "Cart Item", "Update", Constraints{NotNullViolationCode: "quantity"})
	}

	return nil
}

func (r *cartItemRepo) CheckExistence(ctx *gin.Context, db Querier, id int32) error {
	query := `
		SELECT 1 AS exist FROM cart_items
		WHERE id = $1
	`
	var exist int32
	err := db.QueryRow(ctx, query, id).Scan(&exist)
	if err != nil {
		return Parse(err, "Cart Item", "CheckExistence", make(Constraints))
	}
	return nil
}

func (r *cartItemRepo) Delete(ctx *gin.Context, db Querier, id int32) error {
	query := `
		DELETE FROM cart_items
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)
	if err != nil {
		return Parse(err, "Cart Item", "Delete", make(Constraints))
	}
	return nil
}

func (r *cartItemRepo) DeleteByCartID(ctx *gin.Context, db Querier, cartID int32) error {
	query := `
		DELETE FROM cart_items
		WHERE cart_id = $1
	`
	_, err := db.Exec(ctx, query, cartID)
	if err != nil {
		return Parse(err, "Cart Item", "DeleteByCartID", make(Constraints))
	}
	return nil
}
