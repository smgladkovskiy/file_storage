package save

import (
	"context"
	"errors"
	"testing"

	"github.com/smgladkovskiy/file_storage/test/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCommandHandler(t *testing.T) {
	repo := mocks.NewFileStorageRepository(gomock.NewController(t))
	h := NewCommandHandler(repo)

	assert.Equal(t, &CommandHandler{repo: repo}, h)

	require.Panics(t, func() {
		NewCommandHandler(nil)
	})
}

func TestCommandHandler_Run(t *testing.T) {
	t.Parallel()

	repo := mocks.NewFileStorageRepository(gomock.NewController(t))
	h := NewCommandHandler(repo)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand([]byte{1, 2, 3, 4, 5}, "test.txt")
		require.NoError(t, err)

		repo.EXPECT().SaveFile(gomock.Any(), cmd.file).Times(1).Return(nil)

		assert.NoError(t, h.Run(context.TODO(), cmd))
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		badCmd, err := NewCommand([]byte{}, "bad.txt")
		require.NoError(t, err)

		someErr := errors.New("some Err")

		repo.EXPECT().SaveFile(gomock.Any(), badCmd.file).Times(1).Return(someErr)
		assert.ErrorIs(t, someErr, h.Run(context.TODO(), badCmd))
	})
}
