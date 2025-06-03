package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/refine-software/afrad-api/config"
)

const (
	MaxAge = 86400 * 30
)

func InitOauth(env *config.Env) {
	googleClientID := env.GoogleClientID
	googleClientSecret := env.GoogleClientSecret

	store := sessions.NewCookieStore([]byte(env.SessionKey))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = env.Environment == "prod"

	gothic.Store = store

	goth.UseProviders(
		google.New(
			googleClientID,
			googleClientSecret,
			"http://localhost:8080/oauth/google/callback",
			"openid",
			"profile",
			"email",
			"https://www.googleapis.com/auth/user.phonenumbers.read",
		),
	)
}
