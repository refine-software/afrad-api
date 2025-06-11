package database

type SizeRepository interface{}

type sizeRepo struct{}

func NewSizeRepository() SizeRepository {
	return &sizeRepo{}
}
