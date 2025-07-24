package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type sizeReq struct {
	Size  string `json:"size"  binding:"required"`
	Label string `json:"label" binding:"required"`
}

func (s *Server) createSize(ctx *gin.Context) {
	var req sizeReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	sizeRepo := s.DB.Size()

	err = sizeRepo.Create(ctx, db, &models.Size{Size: req.Size, Label: req.Label})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Created(ctx, "size created")
}

func (s *Server) GetSizes(ctx *gin.Context) {
	db := s.DB.Pool()
	sizeRepo := s.DB.Size()

	expectedValues := map[string]bool{
		"فوقي": true,
		"سفلي": true,
		"حذاء": true,
	}

	labels := ctx.QueryArray("label")
	if len(labels) == 0 {
		sizes, err := sizeRepo.GetAll(ctx, db)
		if err != nil {
			apiErr := utils.MapDBErrorToAPIError(err)
			utils.Fail(ctx, apiErr, err)
			return
		}

		utils.Success(ctx, sizes)

	}
	for _, label := range labels {
		if !expectedValues[label] {
			utils.Fail(ctx, utils.NewAPIError(http.StatusBadRequest, "there's no such label"), nil)
			return
		}
	}

	sizes, err := sizeRepo.GetByLabel(ctx, db, labels)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, sizes)
}

type updateSizeReq struct {
	Size  string `json:"size"`
	Label string `json:"label"`
}

func (s *Server) updateSize(ctx *gin.Context) {
	var req updateSizeReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	id := convStrToInt(ctx, ctx.Param("id"), "size_id")

	if req.Size == "" && req.Label == "" {
		utils.Fail(
			ctx,
			utils.NewAPIError(http.StatusBadRequest, "size and label could not be both empty"),
			nil,
		)
		return
	}

	db := s.DB.Pool()
	sizeRepo := s.DB.Size()

	err = sizeRepo.Update(ctx, db, &models.Size{ID: int32(id), Size: req.Size, Label: req.Label})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}

func (s *Server) deleteSize(ctx *gin.Context) {
	id := convStrToInt(ctx, ctx.Param("id"), "size_id")
	db := s.DB.Pool()
	sizeRepo := s.DB.Size()

	err := sizeRepo.Delete(ctx, db, int32(id))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, nil)
}
