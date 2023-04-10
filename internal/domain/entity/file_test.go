package entity_test

import (
	"testing"

	"github.com/smgladkovskiy/file_storage/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFile(t *testing.T) {
	type inStruct struct {
		name  string
		parts int
	}

	type expStruct struct {
		file *entity.File
		err  error
	}

	type testCase struct {
		name string
		in   inStruct
		exp  expStruct
	}

	testCases := []testCase{
		{
			name: "success creating",
			in: inStruct{
				name: "file name",
			},
			exp: expStruct{
				file: &entity.File{
					FileMeta: entity.FileMeta{
						Name: "file name",
					},
				},
				err: nil,
			},
		},
		{
			name: "empty name error",
			in: inStruct{
				name: "",
			},
			exp: expStruct{
				file: nil,
				err:  entity.ErrEmptyFileName,
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := entity.NewFile(tc.in.name)
			if tc.exp.err != nil {
				require.ErrorIs(t, err, tc.exp.err)
			} else {
				require.NoError(t, err)
			}

			if tc.exp.file != nil {
				tc.exp.file.ID = out.ID
			}

			assert.Equal(t, tc.exp.file, out)
		})
	}
}

func TestNewFileFromID(t *testing.T) {
	type expStruct struct {
		file *entity.File
		err  error
	}

	type testCase struct {
		name string
		in   uuid.UUID
		exp  expStruct
	}
	id := uuid.New()
	testCases := []testCase{
		{
			name: "success create",
			in:   id,
			exp: expStruct{
				file: &entity.File{ID: id},
				err:  nil,
			},
		},
		{
			name: "nil id error",
			in:   uuid.UUID{},
			exp: expStruct{
				file: nil,
				err:  entity.ErrNilFileID,
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := entity.NewFileFromID(tc.in)
			if tc.exp.err != nil {
				require.ErrorIs(t, err, tc.exp.err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp.file, out)
		})
	}
}

func TestFile_Combine(t *testing.T) {

	t.Parallel()

	file, err := entity.NewFile("file.name")
	require.NoError(t, err)

	file.SplitContent = [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}, {13, 14, 15}}

	require.NoError(t, file.Combine())

	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, file.Content)
}

func TestFile_SetContent(t *testing.T) {

}

func TestFile_SetSplitContent(t *testing.T) {

}

func TestFile_Split(t *testing.T) {

}
