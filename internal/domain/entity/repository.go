package entity

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=../../../test/mocks/file_storage_repository_mock.go -package=mocks -mock_names FileStorageRepository=FileStorageRepository

type FileStorageRepository interface {
	SaveFile(ctx context.Context, file *File) error
	GetFileByID(ctx context.Context, fileID uuid.UUID) (*File, error)
}
