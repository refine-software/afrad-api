package server

import (
	"errors"
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
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, rr)
}

type updateReviewReq struct {
	Rating int    `json:"rating"`
	Review string `json:"review"`
}

func (s *Server) updateReview(c *gin.Context) {
	reviewID := convStrToInt(c, c.Param("id"), "id")
	if reviewID == 0 {
		return
	}
	var req updateReviewReq
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
	reviewRepo := s.DB.RatingReview()

	rr, err := reviewRepo.Get(c, db, int32(reviewID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	if userID != int(rr.UserID) {
		utils.Fail(c, utils.ErrForbidden, errors.New("a user trying to update other user review"))
		return
	}

	if req.Rating != 0 {
		rr.Rating = int32(req.Rating)
	}

	rr.Review.String = req.Review
	rr.Review.Valid = req.Review != ""

	err = reviewRepo.Update(c, db, rr)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, rr)
}

func (s *Server) deleteReview(c *gin.Context) {
	reviewID := convStrToInt(c, c.Param("id"), "id")
	if reviewID == 0 {
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
	reviewRepo := s.DB.RatingReview()

	rr, err := reviewRepo.Get(c, db, int32(reviewID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	if rr.UserID != int32(userID) {
		utils.Fail(c, utils.ErrForbidden, err)
		return
	}

	err = reviewRepo.Delete(c, db, int32(reviewID))
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, "review deleted successfully")
}
