package database

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/refine-software/afrad-api/config"
)

// Service represents a service that interacts with a database.
type Service interface {
	brand() BrandRepository
	category() CategoryRepository
	color() ColorRepository
	image() ImageRepository
	product() ProductRepository
	productVariant() ProductVariantRepository
	ratingReview() RatingReviewRepository
	size() SizeRepository
	PasswordReset() PasswordResetRepository
	AccountVerificationCode() AccountVerificationCodeRepository
	LocalAuth() LocalAuthRepository
	Oauth() OAuthRepository
	User() UserRepository
	Session() SessionRepository
	Pool() *pgxpool.Pool
	WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error
	BeginTx(ctx *gin.Context) (pgx.Tx, error)
	Close()
}

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type service struct {
	brandRepo                   BrandRepository
	categoryRepo                CategoryRepository
	colorRepo                   ColorRepository
	imageRepo                   ImageRepository
	productRepo                 ProductRepository
	productVariantRepo          ProductVariantRepository
	ratingReviewRepo            RatingReviewRepository
	sizeRepo                    SizeRepository
	accountVerificationCodeRepo AccountVerificationCodeRepository
	passwordResetRepo           PasswordResetRepository
	sessionRepo                 SessionRepository
	localAuthRepo               LocalAuthRepository
	oAuthRepo                   OAuthRepository
	userRepo                    UserRepository
	db                          *pgxpool.Pool
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
		db:                          pool,
		userRepo:                    NewUserRepository(),
		oAuthRepo:                   NewOAuthRepository(),
		localAuthRepo:               NewLocalAuthRepository(),
		sessionRepo:                 NewSessionRepository(),
		passwordResetRepo:           NewPasswordResetRepository(),
		accountVerificationCodeRepo: NewAccountVerificationCodeRepository(),
		brandRepo:                   NewBrandRepository(),
		categoryRepo:                NewCategoryRepository(),
		colorRepo:                   NewColorRepository(),
		imageRepo:                   NewImageRepository(),
		productRepo:                 NewProductRepository(),
		productVariantRepo:          NewProductVariantRepository(),
		ratingReviewRepo:            NewRatingReviewRepository(),
		sizeRepo:                    NewSizeRepository(),
	}

	return dbInstance
}

func (s *service) User() UserRepository {
	return s.userRepo
}

func (s *service) Oauth() OAuthRepository {
	return s.oAuthRepo
}

func (s *service) LocalAuth() LocalAuthRepository {
	return s.localAuthRepo
}

func (s *service) Session() SessionRepository {
	return s.sessionRepo
}

func (s *service) PasswordReset() PasswordResetRepository {
	return s.passwordResetRepo
}

func (s *service) AccountVerificationCode() AccountVerificationCodeRepository {
	return s.accountVerificationCodeRepo
}

func (s *service) brand() BrandRepository {
	return s.brandRepo
}

func (s *service) category() CategoryRepository {
	return s.categoryRepo
}

func (s *service) color() ColorRepository {
	return s.colorRepo
}

func (s *service) image() ImageRepository {
	return s.imageRepo
}

func (s *service) product() ProductRepository {
	return s.productRepo
}

func (s *service) productVariant() ProductVariantRepository {
	return s.productRepo
}

func (s *service) ratingReview() RatingReviewRepository {
	return s.ratingReviewRepo
}

func (s *service) size() SizeRepository {
	return s.sizeRepo
}

func (s *service) Pool() *pgxpool.Pool {
	return s.db
}

func (s *service) WithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (s *service) BeginTx(ctx *gin.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}

func (s *service) Close() {
	s.db.Close()
}
