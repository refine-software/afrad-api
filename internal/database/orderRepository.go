package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type OrderRepository interface {
	Create(ctx *gin.Context, db Querier, order *models.Order) (int32, error)
}

type orderRepo struct{}

func NewOrderRepository() OrderRepository {
	return &orderRepo{}
}

func (r *orderRepo) Create(ctx *gin.Context, db Querier, order *models.Order) (int32, error) {
	query := `
		INSERT INTO orders(name, city_id, town, street, address, phone_number, total_price, order_status, user_id)
		VALUES($1, $2, $3, $4, $5, $6, $7, 'order_placed', $8)
		RETURNING id
	`

	var orderID int32
	err := db.QueryRow(
		ctx,
		query,
		order.Name,
		order.CityID,
		order.Town,
		order.Street,
		order.Address,
		order.PhoneNumber,
		order.TotalPrice,
		order.UserID,
	).Scan(&orderID)
	if err != nil {
		return 0, Parse(
			err,
			"Order",
			"Create",
			Constraints{
				NotNullViolationCode:    "all columns",
				ForeignKeyViolationCode: "user_id, city_id",
			},
		)
	}

	return orderID, nil
}

func (r *orderRepo) GetAll(ctx *gin.Context, db Querier) {
}

func (r *orderRepo) GetByUserID(ctx *gin.Context, db Querier) {
}
