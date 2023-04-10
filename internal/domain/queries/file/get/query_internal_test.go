package get

import (
	"context"
	"errors"
	"testing"

	"github.com/smgladkovskiy/file_storage/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueryHandlerHandler(t *testing.T) {
	repo := mocks.NewFileStorageRepository(gomock.NewController(t))
	h := NewQueryHandler(repo)

	assert.Equal(t, &QueryHandler{repo: repo}, h)

	require.Panics(t, func() {
		NewQueryHandler(nil)
	})
}

func TestCommandHandler_Run(t *testing.T) {
	t.Parallel()

	repo := mocks.NewFileStorageRepository(gomock.NewController(t))
	h := NewQueryHandler(repo)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		q, err := NewQuery(uuid.New())
		require.NoError(t, err)

		file := q.file
		file.SplitContent = [][]byte{{1}, {2}, {3}, {4}, {5}}
		file.Content = []byte{1, 2, 3, 4, 5}
		file.Parts = 5

		repo.EXPECT().GetFileByID(gomock.Any(), q.file.ID).Times(1).Return(file, nil)

		out, err := h.Run(context.TODO(), q)
		assert.NoError(t, err)
		assert.Equal(t, file, out)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		q, err := NewQuery(uuid.New())
		require.NoError(t, err)

		someErr := errors.New("some Error")

		repo.EXPECT().GetFileByID(gomock.Any(), q.file.ID).Times(1).Return(nil, someErr)

		out, err := h.Run(context.TODO(), q)
		assert.ErrorIs(t, someErr, err)
		assert.Nil(t, out)
	})
}
