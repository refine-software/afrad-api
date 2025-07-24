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

type cartReq struct {
	ProductID int32 `json:"productId" binding:"required"`
	Quantity  int   `json:"quantity"  binding:"required"`
}

func (s *Server) addToCart(ctx *gin.Context) {
	var req cartReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	cartRepo := s.DB.Cart()
	cartItemRepo := s.DB.CartItem()
	productVariantRepo := s.DB.ProductVariant()

	claims := auth.GetAccessClaims(ctx)
	if claims == nil {
		return
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	err = s.DB.WithTransaction(ctx, func(tx pgx.Tx) error {
		var (
			cartID        int32
			price         int
			totalPrice    int
			totalQuantity int
		)

		// check if cart exists
		cartID, err = cartRepo.GetIDByUserID(ctx, tx, int32(userID))
		if err != nil && database.IsDBNotFoundErr(err) {
			// if not exists make one
			cartID, err = cartRepo.Create(ctx, tx, int32(userID))
			if err != nil {
				return err
			}
		}

		// calculat total price of one product per quantity
		price, err = productVariantRepo.GetPriceByID(ctx, tx, req.ProductID)
		if err != nil {
			return err
		}
		totalPricePerProduct := price * req.Quantity

		// add product to cart item
		err = cartItemRepo.Create(ctx, tx, &models.CartItem{
			CartID:     cartID,
			ProductID:  req.ProductID,
			Quantity:   req.Quantity,
			TotalPrice: totalPricePerProduct,
		})
		if err != nil {
			return err
		}

		// when we create the cart the (total price) and (quantity) will be both 0.
		// after adding to the cart item we calculat the (total price) and (quantity) of the whole cart
		totalPrice, totalQuantity, err = cartItemRepo.GetPriceQuantityByCartID(ctx, tx, cartID)
		if err != nil {
			return err
		}

		cartRepo.Update(
			ctx,
			tx,
			&models.Cart{
				UserID:     int32(userID),
				TotalPrice: totalPrice,
				Quantity:   totalQuantity,
			})

		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Created(ctx, nil)
}

type cartResponse struct {
	Cart      models.Cart             `json:"cart"`
	CartItems []database.GetCartItems `json:"cartItems"`
}

func (s *Server) getCart(ctx *gin.Context) {
	var res cartResponse

	db := s.DB.Pool()
	cartRepo := s.DB.Cart()
	cartItemRepo := s.DB.CartItem()

	claims := auth.GetAccessClaims(ctx)
	if claims == nil {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}
	fmt.Println(userID)

	cart, err := cartRepo.GetByUserID(ctx, db, int32(userID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	cartItems, err := cartItemRepo.GetAll(ctx, db, cart.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	res = cartResponse{
		Cart:      *cart,
		CartItems: cartItems,
	}

	utils.Success(ctx, res)
}

type updateQantityReq struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (s *Server) updateCartItemQuantity(ctx *gin.Context) {
	var req updateQantityReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	id := convStrToInt(ctx, ctx.Param("id"), "cart_item_id")

	db := s.DB.Pool()
	cartItemRepo := s.DB.CartItem()

	err = cartItemRepo.Update(ctx, db, int32(id), req.Quantity)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, nil)
}

func (s *Server) deleteCartItem(ctx *gin.Context) {
	id := convStrToInt(ctx, ctx.Param("id"), "cart item id")
	if id == 0 {
		return
	}

	db := s.DB.Pool()
	cartItemRepo := s.DB.CartItem()

	err := cartItemRepo.Delete(ctx, db, int32(id))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, nil)
}

func (s *Server) deleteCart(ctx *gin.Context) {
	claims := auth.GetAccessClaims(ctx)
	if claims == nil {
		return
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
		return
	}

	cartRepo := s.DB.Cart()
	db := s.DB.Pool()

	err = cartRepo.Delete(ctx, db, int32(userID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}
