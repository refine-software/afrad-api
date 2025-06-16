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
	ParentID *int32 `json:"parentID"`
}

func (s *Server) createCategory(ctx *gin.Context) {
	var req categoryRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(ctx, utils.ErrBadRequest, err)
		return
	}

	categoryRepo := s.db.Category()
	db := s.db.Pool()

	parentID := pgtype.Int4{Valid: false}
	if req.ParentID != nil {
		parentID = pgtype.Int4{Int32: *req.ParentID, Valid: true}
	}
	c := &models.Category{
		Name:     req.Name,
		ParentID: parentID,
	}

	_, dbErr := categoryRepo.Create(ctx, db, c)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "category")
		utils.Fail(ctx, apiErr, dbErr)
		return
	}
	utils.Created(ctx, "category created")
}

type getCategoriesRes struct {
	Categories []models.Category `json:"categories"`
}

func (s *Server) getCategories(ctx *gin.Context) {
	categoryRepo := s.db.Category()
	db := s.db.Pool()

	categories, dbErr := categoryRepo.GetAll(ctx, db)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "categories")
		utils.Fail(ctx, apiErr, dbErr)
		return
	}

	utils.Success(ctx, getCategoriesRes{
		Categories: *categories,
	})
}

func (s *Server) deleteCategory(ctx *gin.Context) {
	categoryRepo := s.db.Category()
	db := s.db.Pool()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
	}

	dbErr := categoryRepo.CheckExistence(ctx, db, int32(id))
	if dbErr != nil && dbErr.Message == database.ErrNotFound {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "category doesn't exists",
			},
			dbErr,
		)
	}

	dbErr = categoryRepo.Delete(ctx, db, int32(id))
	if dbErr != nil && dbErr.Message == database.ErrForeignKey {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "can not delete this category, you should delete its childs first",
			},
			dbErr,
		)
	}

	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "category")
		utils.Fail(ctx, apiErr, dbErr)
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

	categoryRepo := s.db.Category()
	db := s.db.Pool()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(ctx, utils.ErrInternal, err)
	}

	dbErr := categoryRepo.CheckExistence(ctx, db, int32(id))
	if dbErr != nil && dbErr.Message == database.ErrNotFound {
		utils.Fail(
			ctx,
			&utils.APIError{
				Code:    http.StatusBadRequest,
				Message: "category doesn't exists",
			},
			dbErr,
		)
		return
	}

	dbErr = categoryRepo.Update(ctx, db, int32(id), req.Name)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "category")
		utils.Fail(ctx, apiErr, dbErr)
		return
	}

	utils.Success(ctx, nil)
}
