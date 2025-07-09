// Package database provides access to the application's PostgreSQL database and
// defines the service layer for interacting with various data repositories.
// It offers a singleton service instance that exposes repository interfaces for
// different domain models such as users, products, authentication, and more.
//
// This package also handles connection pooling, transaction management, and
// context-aware query execution using the pgx PostgreSQL driver.
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
	City() CityRepository
	OrderDetails() OrderDetailsRepository
	Order() OrderRepository
	CartItem() CartItemRepository
	Cart() CartRepository
	Brand() BrandRepository
	Category() CategoryRepository
	Color() ColorRepository
	Image() ImageRepository
	Product() ProductRepository
	ProductVariant() ProductVariantRepository
	RatingReview() RatingReviewRepository
	Size() SizeRepository
	PasswordReset() PasswordResetRepository
	AccountVerificationCode() AccountVerificationCodeRepository
	LocalAuth() LocalAuthRepository
	Oauth() OAuthRepository
	User() UserRepository
	Session() SessionRepository
	Wishlist() WishlistRepository
	Pool() *pgxpool.Pool
	// Make sure to use this method when all errors being returned are db errors.
	// you can use it when other errors are being returned but still.
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
	cityRepo                    CityRepository
	orderDetailsRepo            OrderDetailsRepository
	orderRepo                   OrderRepository
	cartItemRepo                CartItemRepository
	cartRepo                    CartRepository
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
	wishlistRepo                WishlistRepository
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
		wishlistRepo:                NewWishlistRepository(),
		cartRepo:                    NewCartRepository(),
		cartItemRepo:                NewCartItemRepository(),
		orderRepo:                   NewOrderRepository(),
		orderDetailsRepo:            NewOrderDetailsRepository(),
		cityRepo:                    NewCityRepository(),
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

func (s *service) Brand() BrandRepository {
	return s.brandRepo
}

func (s *service) Category() CategoryRepository {
	return s.categoryRepo
}

func (s *service) Color() ColorRepository {
	return s.colorRepo
}

func (s *service) Image() ImageRepository {
	return s.imageRepo
}

func (s *service) Product() ProductRepository {
	return s.productRepo
}

func (s *service) ProductVariant() ProductVariantRepository {
	return s.productVariantRepo
}

func (s *service) RatingReview() RatingReviewRepository {
	return s.ratingReviewRepo
}

func (s *service) Size() SizeRepository {
	return s.sizeRepo
}

func (s *service) Wishlist() WishlistRepository {
	return s.wishlistRepo
}

func (s *service) Cart() CartRepository {
	return s.cartRepo
}

func (s *service) CartItem() CartItemRepository {
	return s.cartItemRepo
}

func (s *service) Order() OrderRepository {
	return s.orderRepo
}

func (s *service) OrderDetails() OrderDetailsRepository {
	return s.orderDetailsRepo
}

func (s *service) City() CityRepository {
	return s.cityRepo
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
