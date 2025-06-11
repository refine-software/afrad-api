package database

type BrandRepository interface{}

type brandRepo struct{}

func NewBrandRepository() BrandRepository {
	return &brandRepo{}
}
