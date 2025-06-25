package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/database"
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

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	colorRepo := s.DB.Color()

	err = colorRepo.CheckExistenece(ctx, db, int32(id))
	if err != nil && database.IsDBNotFoundErr(err) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "color not found",
			},
			err,
		)
		return
	}

	err = colorRepo.Update(ctx, db, int32(id), req.Color)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}

func (s *Server) deleteColor(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	colorRepo := s.DB.Color()

	err = colorRepo.CheckExistenece(ctx, db, int32(id))
	if err != nil && database.IsDBNotFoundErr(err) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "color not found",
			},
			err,
		)
		return
	}

	err = colorRepo.Delete(ctx, db, int32(id))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Success(ctx, nil)
}
