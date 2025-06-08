package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	// APP
	Environment string `mapstructure:"APP_ENV"`
	Port        int    `mapstructure:"PORT"`

	// DB
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_DATABASE"`
	DBUsername string `mapstructure:"DB_USERNAME"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBSchema   string `mapstructure:"DB_SCHEMA"`
	SSLMode    string `mapstructure:"SSL_MODE"`
	DBUrl      string

	// S3
	S3AccessKey       string `mapstructure:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey string `mapstructure:"S3_SECRET_ACCESS_KEY"`
	S3Region          string `mapstructure:"S3_REGION"`
	S3Bucket          string `mapstructure:"S3_BUCKET"`

	// Oauth
	SessionKey         string `mapstructure:"SESSION_KEY"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`

	// JWT
	AccessTokenSecret     string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret    string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenExpInMin   int    `mapstructure:"ACCESS_TOKEN_EXP_IN_MIN"`
	RefreshTokenExpInDays int    `mapstructure:"REFRESH_TOKEN_EXP_IN_DAYS"`

	// Hash
	HashSecret string `mapstructure:"HASHING_SECRET"`

	// Email
	Email    string `mapstructure:"EMAIL"`
	Password string `mapstructure:"PASSWORD"`
	AppURL   string `mapstructure:"APP_URL"`

	// Twilio
	TwilioAccountSID  string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken   string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioSandboxFrom string `mapstructure:"TWILIO_SANDBOX_FROM"`
}

func NewEnv() *Env {
	env := Env{}
	viper.AutomaticEnv()
	viper.SetDefault("APP_ENV", "dev")

	bindEnvVariables()

	err := viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded:", err)
	}

	switch env.Environment {
	case "dev":
		fmt.Println("You're running your application in development mode")
	case "prod":
		fmt.Println("You're running your application in production mode")
	default:
		log.Fatalf(
			"Unknown environment: %s. Please set APP_ENV to 'dev' or 'prod'",
			env.Environment,
		)
	}

	env.DBUrl = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		env.DBUsername,
		env.DBPassword,
		env.DBHost,
		env.DBPort,
		env.DBName,
		env.SSLMode,
	)

	if env.Environment == "dev" {
		log.Println("The App is running in development mode")
	}

	return &env
}

// viper.AutomaticEnv() only automatically picks up environment variables if their keys match exactly —
// and by default, Viper expects the struct field names, not the mapstructure tags.
// But Viper doesn't automatically know to use the mapstructure keys for env lookup.
func bindEnvVariables() {
	vars := []string{
		"APP_ENV",
		"PORT",
		// DB
		"DB_HOST",
		"DB_PORT",
		"DB_DATABASE",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_SCHEMA",
		"SSL_MODE",
		// S3
		"S3_ACCESS_KEY_ID",
		"S3_SECRET_ACCESS_KEY",
		"S3_REGION",
		"S3_BUCKET",
		// Oauth
		"SESSION_KEY",
		"GOOGLE_CLIENT_ID",
		"GOOGLE_CLIENT_SECRET",
		// JWT
		"ACCESS_TOKEN_SECRET",
		"REFRESH_TOKEN_SECRET",
		"ACCESS_TOKEN_EXP_IN_MIN",
		"REFRESH_TOKEN_EXP_IN_DAYS",
		// Hashing
		"HASHING_SECRET",
		// Twilio
		"TWILIO_ACCOUNT_SID",
		"TWILIO_AUTH_TOKEN",
		"TWILIO_SANDBOX_FROM",
		// email
		"EMAIL",
		"PASSWORD",
		"APP_URL",
	}

	for _, key := range vars {
		if err := viper.BindEnv(key); err != nil {
			log.Fatalf("Failed to bind environment variable %s: %v", key, err)
		}
	}
}
