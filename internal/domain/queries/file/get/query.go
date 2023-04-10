package get

import (
	"github.com/smgladkovskiy/file_storage/internal/domain/entity"

	"github.com/google/uuid"
)

type Query struct {
	file *entity.File
}

func NewQuery(id uuid.UUID) (*Query, error) {
	file, err := entity.NewFileFromID(id)
	if err != nil {
		return nil, err
	}

	return &Query{file: file}, nil
}
