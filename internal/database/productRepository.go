package database

type ProductRepository interface{}

type productRepo struct{}

func NewProductRepository() ProductRepository {
	return &productRepo{}
}
