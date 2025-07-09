package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type OrderDetailsRepository interface {
	// Create order details
	//
	// required columns: product_id, quantity, total_price, order_id
	Create(ctx *gin.Context, db Querier, orderDetails *models.OrderDetails) error
}

type orderDetailsRepo struct{}

func NewOrderDetailsRepository() OrderDetailsRepository {
	return &orderDetailsRepo{}
}

func (r *orderDetailsRepo) Create(
	ctx *gin.Context,
	db Querier,
	orderDetails *models.OrderDetails,
) error {
	query := `
		INSERT INTO order_details(product_id, quantity, total_price, order_id)
		VALUES($1, $2, $3, $4)
	`

	_, err := db.Exec(
		ctx,
		query,
		orderDetails.ProductID,
		orderDetails.Quantity,
		orderDetails.TotalPrice,
		orderDetails.OrderID,
	)
	if err != nil {
		return Parse(
			err,
			"orderDetails",
			"Create",
			Constraints{
				NotNullViolationCode:    "all columns",
				UniqueViolationCode:     "order_id & product_id",
				ForeignKeyViolationCode: "order_id, product_id",
			},
		)
	}
	return nil
}
