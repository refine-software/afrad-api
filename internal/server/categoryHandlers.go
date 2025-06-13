package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

type getCategoryRes struct {
	Categories []models.Category `json:"categories"`
}

func (s *Server) getCategories(ctx *gin.Context) {
	categoryRepo := s.db.Category()
	db := s.db.Pool()

	categories, dbErr := categoryRepo.Get(ctx, db)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "categories")
		utils.Fail(ctx, apiErr, dbErr)
		return
	}

	utils.Success(ctx, getCategoryRes{
		Categories: *categories,
	})
}
