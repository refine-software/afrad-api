package database

type ImageRepository interface{}

type imageRepo struct{}

func NewImageRepository() ImageRepository {
	return &imageRepo{}
}
