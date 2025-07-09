package server

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type variantReq struct {
	Quantity  int   `json:"quantity"  binding:"required"`
	Price     int   `json:"price"     binding:"required"`
	ColorID   int32 `json:"colorId"   binding:"required"`
	SizeID    int32 `json:"sizeId"    binding:"required"`
	ProductID int32 `json:"productId" binding:"required"`
}

func (s *Server) addVariant(c *gin.Context) {
	var req variantReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	variantRepo := s.DB.ProductVariant()

	pv := models.ProductVariant{
		Quantity:  req.Quantity,
		Price:     req.Price,
		ColorID:   req.ColorID,
		SizeID:    req.SizeID,
		ProductID: req.ProductID,
	}

	err = variantRepo.Create(c, db, &pv)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, "variant added successfully")
}

func (s *Server) getVariant(c *gin.Context) {
	variantID := convStrToInt(c, c.Param("id"), "variant id")
	if variantID == 0 {
		return
	}

	db := s.DB.Pool()
	variantRepo := s.DB.ProductVariant()

	pv, err := variantRepo.Get(c, db, int32(variantID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, pv)
}

func (s *Server) deleteVariant(c *gin.Context) {
	variantID := convStrToInt(c, c.Param("id"), "variant id")
	if variantID == 0 {
		return
	}

	db := s.DB.Pool()
	variantRepo := s.DB.ProductVariant()

	err := variantRepo.Delete(c, db, variantID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.NoContent(c)
}

type productVariantReq struct {
	Quantity int   `json:"quantity"`
	Price    int   `json:"price"`
	ColorID  int32 `json:"colorId"`
	SizeID   int32 `json:"sizeId"`
}

func (s *Server) updateVariant(c *gin.Context) {
	variantID := int32(convStrToInt(c, c.Param("id"), "variant id"))
	if variantID == 0 {
		return
	}

	var req productVariantReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	variantRepo := s.DB.ProductVariant()

	pv := models.ProductVariant{
		ID:       variantID,
		Quantity: req.Quantity,
		Price:    req.Price,
		ColorID:  req.ColorID,
		SizeID:   req.SizeID,
	}

	err = variantRepo.Update(c, db, &pv)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, "updated successfully")
}
