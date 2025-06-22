package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type OAuthRepository interface {
	Create(ctx *gin.Context, db Querier, oauth *models.OAuth) error
}

type oAuthRepo struct{}

func NewOAuthRepository() OAuthRepository {
	return &oAuthRepo{}
}

func (a *oAuthRepo) Create(
	ctx *gin.Context,
	db Querier,
	oauth *models.OAuth,
) error {
	query := `
		INSERT INTO oauth(user_id, provider, provider_id)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(ctx, query, oauth.UserID, oauth.Provider, oauth.ProviderID)
	if err != nil {
		return Parse(err, "OAuth", "Create", Constraints{
			UniqueViolationCode:     "provider_id", // unique index on provider_id
			ForeignKeyViolationCode: "user",        // FK to users
			NotNullViolationCode:    "provider",
			UniqueViolationCode:     "user", // primary key on user_id acts like unique constraint
		})
	}
	return nil
}
