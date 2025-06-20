package server

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/utils"
)

func (s *Server) getUser(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	db := s.DB.Pool()
	userRepo := s.DB.User()

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	user, err := userRepo.Get(c, db, userID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, user)
}

type updateUserReq struct {
	FirstName string `form:"firstName"`
	LastName  string `form:"lastName"`
}

func (s *Server) updateUser(c *gin.Context) {
	var req updateUserReq
	err := c.ShouldBind(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	db := s.DB.Pool()
	userRepo := s.DB.User()
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	user, dbErr := userRepo.Get(c, db, userID)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "user")
		utils.Fail(c, apiErr, dbErr)
		return
	}

	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = pgtype.Text{String: req.LastName, Valid: true}
	}

	imageUpload, apiErr := getImageFile(c)
	if apiErr != nil {
		utils.Fail(c, apiErr, errors.New(apiErr.Message))
		return
	}

	err = s.DB.WithTransaction(c, func(tx pgx.Tx) error {
		// if image exists do the following:
		// delete old image if exists
		if imageUpload != nil {
			if user.Image.Valid {
				_ = s.S3.DeleteImageByURL(c, user.Image.String)
			}
			// upload new image
			var newImageURL string
			newImageURL, err = s.S3.UploadImage(c, imageUpload.File, imageUpload.Header)
			if err != nil {
				return err
			}

			// update image in user struct
			user.Image.String = newImageURL
			user.Image.Valid = true
		}

		dbErr = userRepo.Update(c, db, user)
		if dbErr != nil {
			return dbErr
		}

		return nil
	})
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "user")
		utils.Fail(c, apiErr, dbErr)
		return
	}

	utils.Success(c, user)
}

func (s *Server) deleteUser(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	db := s.DB.Pool()
	userRepo := s.DB.User()
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	u, err := userRepo.Get(c, db, userID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	err = userRepo.Delete(c, db, userID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err, "user")
		utils.Fail(c, apiErr, err)
		return
	}

	if u.Image.Valid {
		_ = s.S3.DeleteImageByURL(c, u.Image.String)
	}

	utils.Success(c, nil)
}
