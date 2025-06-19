package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/auth"
	"github.com/refine-software/afrad-api/internal/utils"
)

func (s *Server) getUser(c *gin.Context) {
	claims := auth.GetAccessClaims(c)
	if claims == nil {
		return
	}

	db := s.db.Pool()
	userRepo := s.db.User()

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
