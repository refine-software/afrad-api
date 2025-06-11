package database

type CategoryRepository interface{}

type categoryRepo struct{}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepo{}
}
