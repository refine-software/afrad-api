package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/models"
	"github.com/refine-software/afrad-api/internal/utils"
)

type ReviewReq struct {
	Rating    int    `json:"rating"    binding:"required"`
	Review    string `json:"review"`
	ProductID int    `json:"productId" binding:"required"`
}

func (s *Server) postReview(c *gin.Context) {
	var req ReviewReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

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
	ratingReviewRepo := s.DB.RatingReview()

	rr := models.RatingReview{
		Rating:    int32(req.Rating),
		Review:    pgtype.Text{String: req.Review, Valid: req.Review != ""},
		ProductID: int32(req.ProductID),
		UserID:    int32(userID),
	}

	err = ratingReviewRepo.Create(c, db, &rr)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "product")
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, rr)
}
