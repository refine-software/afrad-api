package server

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/utils"
)

type colorReq struct {
	Color string `json:"color" binding:"required"`
}

func (s *Server) createColor(ctx *gin.Context) {
	var req colorReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	colorRepo := s.DB.Color()

	err = colorRepo.Create(ctx, db, req.Color)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Created(ctx, nil)
}

func (s *Server) getColors(ctx *gin.Context) {
	db := s.DB.Pool()
	colorRepo := s.DB.Color()

	colors, err := colorRepo.GetAll(ctx, db)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, colors)
}

func (s *Server) updateColor(ctx *gin.Context) {
	var req colorReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	id := convStrToInt(ctx, ctx.Param("id"), "color_id")

	db := s.DB.Pool()
	colorRepo := s.DB.Color()

	err = colorRepo.Update(ctx, db, int32(id), req.Color)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}

func (s *Server) deleteColor(ctx *gin.Context) {
	id := convStrToInt(ctx, ctx.Param("id"), "color_id")
	db := s.DB.Pool()
	colorRepo := s.DB.Color()

	err := colorRepo.Delete(ctx, db, int32(id))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}
