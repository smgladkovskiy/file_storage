package get

import (
	"context"

	"github.com/smgladkovskiy/file_storage/internal/domain/entity"
)

type QueryHandler struct {
	repo entity.FileStorageRepository
}

func NewQueryHandler(repo entity.FileStorageRepository) *QueryHandler {
	if repo == nil {
		panic("empty FileStorageRepository")
	}

	return &QueryHandler{repo: repo}
}

func (h *QueryHandler) Run(ctx context.Context, query *Query) (*entity.File, error) {
	return h.repo.GetFileByID(ctx, query.file.ID)
}
