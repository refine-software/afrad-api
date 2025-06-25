package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

func (s *Server) getWishlist(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	db := s.DB.Pool()
	wishlistRepo := s.DB.Wishlist()
	ws, err := wishlistRepo.GetAllOfUser(c, db, int32(userID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	if len(ws) == 0 {
		utils.NoContent(c)
		return
	}

	utils.Success(c, ws)
}

func (s *Server) addToWishlist(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	productID := convStrToInt(c, c.Param("product_id"), "product_id")
	if productID == 0 {
		return
	}

	db := s.DB.Pool()
	wishlistRepo := s.DB.Wishlist()

	w := models.Wishlist{
		UserID:    int32(userID),
		ProductID: int32(productID),
	}

	err = wishlistRepo.Create(c, db, &w)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, "product added to wishlist")
}

func (s *Server) deleteFromWishlist(c *gin.Context) {
	id := convStrToInt(c, c.Param("id"), "id")
	if id == 0 {
		return
	}

	db := s.DB.Pool()
	wishlistRepo := s.DB.Wishlist()

	err := wishlistRepo.Delete(c, db, int32(id))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, "wishlist item is deleted successfully")
}
