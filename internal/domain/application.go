package domain

import (
	"github.com/smgladkovskiy/file_storage/internal/domain/commands/file/save"
	"github.com/smgladkovskiy/file_storage/internal/domain/queries/file/get"
	"github.com/smgladkovskiy/file_storage/internal/domain/repository/inmemory"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SaveFile *save.CommandHandler
}

type Queries struct {
	GetFile *get.QueryHandler
}

func NewApplication(storagesAmount int) *Application {
	fileRepo := inmemory.NewInMemoryStorage(storagesAmount)

	appl := &Application{
		Commands: Commands{
			SaveFile: save.NewCommandHandler(fileRepo),
		},

		Queries: Queries{
			GetFile: get.NewQueryHandler(fileRepo),
		},
	}

	return appl
}
