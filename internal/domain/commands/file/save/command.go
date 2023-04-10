package save

import (
	"github.com/smgladkovskiy/file_storage/internal/domain/entity"

	"github.com/google/uuid"
)

type Command struct {
	file *entity.File
}

func NewCommand(fileContent []byte, fileName string) (*Command, error) {
	file, err := entity.NewFile(fileName)
	if err != nil {
		return nil, err
	}

	file.SetContent(fileContent)

	return &Command{file: file}, nil
}

func (c Command) GetFileID() uuid.UUID {
	return c.file.ID
}
