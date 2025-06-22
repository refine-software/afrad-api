package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type categoryRequest struct {
	Name     string `json:"name"     binding:"required"`
	ParentID int32  `json:"parentId"`
}

func (s *Server) createCategory(ctx *gin.Context) {
	var req categoryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	categoryRepo := s.DB.Category()
	db := s.DB.Pool()

	parentID := pgtype.Int4{Int32: req.ParentID, Valid: req.ParentID != 0}
	c := &models.Category{
		Name:     req.Name,
		ParentID: parentID,
	}

	_, err = categoryRepo.Create(ctx, db, c)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}
	utils.Created(ctx, "category created")
}

type getCategoriesRes struct {
	Categories []models.Category `json:"categories"`
}

func (s *Server) getCategories(ctx *gin.Context) {
	categoryRepo := s.DB.Category()
	db := s.DB.Pool()

	categories, err := categoryRepo.GetAll(ctx, db)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, getCategoriesRes{
		Categories: *categories,
	})
}

func (s *Server) deleteCategory(ctx *gin.Context) {
	categoryRepo := s.DB.Category()
	db := s.DB.Pool()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
	}

	err = categoryRepo.CheckExistence(ctx, db, int32(id))
	if err != nil && database.IsDBNotFoundErr(err) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "category doesn't exists",
			},
			err,
		)
	}

	err = categoryRepo.Delete(ctx, db, int32(id))
	if err != nil && database.IsDBNotFoundErr(err) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "can not delete this category, you should delete its childs first",
			},
			err,
		)
	}

	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, nil)
}

type updateReq struct {
	Name string `json:"name" binding:"required"`
}

func (s *Server) updateCategory(ctx *gin.Context) {
	var req categoryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	categoryRepo := s.DB.Category()
	db := s.DB.Pool()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
	}

	err = categoryRepo.CheckExistence(ctx, db, int32(id))
	if err != nil && database.IsDBNotFoundErr(err) {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "category doesn't exists",
			},
			err,
		)
		return
	}

	err = categoryRepo.Update(ctx, db, int32(id), req.Name)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(ctx, apiErr, err)
		return
	}

	utils.Success(ctx, nil)
}
