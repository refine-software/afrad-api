package database

type ProductVariantRepository interface{}

type productVariantRepo struct{}

func NewProductVariantRepository() ProductVariantRepository {
	return &productVariantRepo{}
}
