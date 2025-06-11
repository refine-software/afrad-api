package database

type ColorRepository interface{}

type colorRepo struct{}

func NewColorRepository() CategoryRepository {
	return &categoryRepo{}
}
