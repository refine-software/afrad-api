package database

import (
	"context"
	"log"

	"afrad-api/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	Close()
}

type service struct {
	db *pgxpool.Pool
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
		db: pool,
	}

	return dbInstance
}

func (s *service) Close() {
	s.db.Close()
}
