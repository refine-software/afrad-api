package database

import (
	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/models"
)

type SessionRepository interface {
	// Get a single session.
	//
	// By: user_id, user_agent.
	GetByUserIDAndUserAgent(
		ctx *gin.Context,
		db Querier,
		userID int32,
		userAgent string,
	) (models.Session, error)

	// Fetch all sessions of a certain user
	GetAllOfUser(
		ctx *gin.Context,
		db Querier,
		userID int32,
	) ([]models.Session, error)

	// This method will create a session, the following columns are required:
	// user_id, user_agent, refresh_token, expires_at.
	Create(ctx *gin.Context, db Querier, sess *models.Session) error

	// Update the user session with the following columns:
	// revoked, refresh_token, expires_at.
	// the session will be updated using its id.
	Update(ctx *gin.Context, db Querier, sess *models.Session) error

	// This method will revoke all sessions of a certain user
	RevokeAllOfUser(ctx *gin.Context, db Querier, userID int32) error
}

type sessionRepo struct{}

func NewSessionRepository() SessionRepository {
	return &sessionRepo{}
}

func (sr *sessionRepo) GetByUserIDAndUserAgent(
	ctx *gin.Context,
	db Querier,
	userID int32,
	userAgent string,
) (models.Session, error) {
	query := `
	SELECT id, revoked, user_agent, refresh_token, expires_at, created_at, updated_at, user_id 
	FROM sessions 
	WHERE user_id = $1 AND user_agent = $2
	`

	var s models.Session
	err := db.QueryRow(ctx, query, userID, userAgent).Scan(
		&s.ID,
		&s.Revoked,
		&s.UserAgent,
		&s.RefreshToken,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.UserID,
	)

	return s, Parse(err, "Session", "GetByUserIDAndUserAgent", make(Constraints))
}

func (sr *sessionRepo) GetAllOfUser(
	ctx *gin.Context,
	db Querier,
	userID int32,
) ([]models.Session, error) {
	query := `
	SELECT id, revoked, user_agent, refresh_token, expires_at, created_at, updated_at, user_id
	FROM sessions
	WHERE user_id = $1
	`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, Parse(err, "Session", "GetAllOfUser", make(Constraints))
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		if err = rows.Scan(&s.ID, &s.Revoked, &s.UserAgent, &s.RefreshToken, &s.ExpiresAt, &s.CreatedAt, &s.UpdatedAt, &s.UserID); err != nil {
			return nil, Parse(err, "Session", "GetAllOfUser", make(Constraints))
		}
		sessions = append(sessions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, Parse(err, "Session", "GetAllOfUser", make(Constraints))
	}

	return sessions, nil
}

func (sr *sessionRepo) Create(ctx *gin.Context, db Querier, sess *models.Session) error {
	query := `
	INSERT INTO sessions(user_id, user_agent, refresh_token, expires_at)
	VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(ctx, query, sess.UserID, sess.UserAgent, sess.RefreshToken, sess.ExpiresAt)
	return Parse(err, "Session", "Create", Constraints{
		UniqueViolationCode:     "user_agent", // because the unique constraint is on (user_id, user_agent)
		ForeignKeyViolationCode: "user_id",
		NotNullViolationCode:    "user_id", // or handle multiple fields if you want
	})
}

func (sr *sessionRepo) Update(ctx *gin.Context, db Querier, sess *models.Session) error {
	query := `
	UPDATE sessions
	SET revoked = $2, refresh_token = $3, expires_at = $4
	WHERE id = $1
	`

	_, err := db.Exec(ctx, query, sess.ID, sess.Revoked, sess.RefreshToken, sess.ExpiresAt)
	return Parse(err, "Session", "Update", Constraints{
		NotNullViolationCode: "refresh_token",
	})
}

func (sr *sessionRepo) RevokeAllOfUser(ctx *gin.Context, db Querier, userID int32) error {
	query := `
		UPDATE sessions
		SET revoked = true
		WHERE user_id = $1
	`

	_, err := db.Exec(ctx, query, userID)
	if err != nil {
		return Parse(err, "Session", "RevokeAllOfUser", make(Constraints))
	}

	return nil
}
