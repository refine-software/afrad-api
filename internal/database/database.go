package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/refine-software/afrad-api/config"
)

// Service represents a service that interacts with a database.
type Service interface {
	AuthProvider() AuthProviderRepository
	User() UserRepository
	Close()
}

type service struct {
	authProviderRepo AuthProviderRepository
	userRepo         UserRepository
	db               *pgxpool.Pool
}

var dbInstance *service

func New(env *config.Env) Service {
	if dbInstance != nil {
		return dbInstance
	}

	pgxConf, err := pgxpool.ParseConfig(env.DBUrl)
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL:", err)
	}

	pgxConf.MaxConns = 10
	pgxConf.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConf)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db:               pool,
		userRepo:         NewUserRepository(pool),
		authProviderRepo: NewAuthProviderRepo(pool),
	}

	return dbInstance
}

func (s *service) User() UserRepository {
	return s.userRepo
}

func (s *service) AuthProvider() AuthProviderRepository {
	return s.authProviderRepo
}

func (s *service) Close() {
	s.db.Close()
}
