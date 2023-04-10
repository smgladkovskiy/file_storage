package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/smgladkovskiy/file_storage/internal/domain"
	"github.com/smgladkovskiy/file_storage/internal/domain/commands/file/save"
	"github.com/smgladkovskiy/file_storage/internal/domain/queries/file/get"
	"github.com/smgladkovskiy/file_storage/pkg/server/openapi"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

type MuxHandlers struct {
	app *domain.Application
}

var (
	ErrShouldBeAttachment = errors.New("content-disposition should be 'form-data")
	ErrNoFileName         = errors.New("no file name in content-disposition")
)

func NewHTTPServer(app *domain.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", openapi.Handler(newHandlers(app)))

	return r
}

func newHandlers(app *domain.Application) *MuxHandlers {
	return &MuxHandlers{app: app}
}

func (srv *MuxHandlers) PutFile(w http.ResponseWriter, r *http.Request) {
	fileContent, fileName, err := readFileFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd, err := save.NewCommand(fileContent, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = srv.app.Commands.SaveFile.Run(r.Context(), cmd); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		FileID uuid.UUID `json:"fileID"`
	}{
		FileID: cmd.GetFileID(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func (srv *MuxHandlers) GetFileFileID(w http.ResponseWriter, r *http.Request, fileID openapi_types.UUID) {
	fileID, err := readFileIDFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	q, err := get.NewQuery(fileID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	file, err := srv.app.Queries.GetFile.Run(r.Context(), q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(file.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

}

func readFileIDFromRequest(r *http.Request) (uuid.UUID, error) {
	fileIDStr := chi.URLParam(r, "fileID")

	return uuid.Parse(fileIDStr)
}

func readFileFromRequest(r *http.Request) ([]byte, string, error) {
	disposition, params, err := mime.ParseMediaType(r.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, "", err
	}

	if disposition != "form-data" {
		return nil, "", fmt.Errorf("%w, %s is given", ErrShouldBeAttachment, disposition)
	}

	filename, ok := params["filename"]
	if !ok {
		return nil, "", ErrNoFileName
	}

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, "", err
	}

	return bb, filename, nil
}
