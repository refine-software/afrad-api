package server

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type orderReq struct {
	Name            string `json:"name"        binding:"required"`
	CityID          int32  `json:"cityId"      binding:"required"`
	Town            string `json:"town"        binding:"required"`
	Street          string `json:"street"      binding:"required"`
	Address         string `json:"address"     binding:"required"`
	PhoneNumber     string `json:"phoneNumber" binding:"required"`
	OrderTotalPrice int    `json:"totalPrice"  binding:"required"`
}

func (s *Server) createOrder(ctx *gin.Context) {
	var req orderReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	claims := auth.GetAccessClaims(ctx)
	if claims == nil {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	orderRepo := s.DB.Order()
	orderDetailsRepo := s.DB.OrderDetails()
	cartRepo := s.DB.Cart()
	cartItemRepo := s.DB.CartItem()

	err = s.DB.WithTransaction(ctx, func(tx pgx.Tx) error {
		var (
			orderID   int32
			cartID    int32
			cartItems []database.GetCartItems
		)
		// create the order
		orderID, err = orderRepo.Create(
			ctx,
			tx,
			&models.Order{
				Name:        req.Name,
				CityID:      req.CityID,
				Town:        req.Town,
				Street:      req.Street,
				Address:     req.Address,
				PhoneNumber: req.PhoneNumber,
				TotalPrice:  req.OrderTotalPrice,
				UserID:      int32(userID),
			},
		)
		if err != nil {
			return err
		}

		// get cart items
		cartID, err = cartRepo.GetIDByUserID(ctx, tx, int32(userID))
		if err != nil {
			return err
		}
		cartItems, err = cartItemRepo.GetAll(ctx, tx, cartID)
		if err != nil {
			return err
		}

		for _, item := range cartItems {
			fmt.Println(item.VariantID)
			err = orderDetailsRepo.Create(
				ctx,
				tx,
				&models.OrderDetails{
					Quantity:   item.Quantity,
					ProductID:  item.VariantID,
					TotalPrice: item.TotalPrice,
					OrderID:    orderID,
				},
			)
			if err != nil {
				return err
			}
		}

		err = cartRepo.Delete(ctx, tx, cartID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}
