package save

import (
	"context"

	"github.com/smgladkovskiy/file_storage/internal/domain/entity"
)

type CommandHandler struct {
	repo entity.FileStorageRepository
}

func NewCommandHandler(repo entity.FileStorageRepository) *CommandHandler {
	if repo == nil {
		panic("empty FileStorageRepository")
	}

	return &CommandHandler{repo: repo}
}

func (h *CommandHandler) Run(ctx context.Context, cmd *Command) error {
	return h.repo.SaveFile(ctx, cmd.file)
}
