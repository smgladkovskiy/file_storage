package file_test

import (
	"testing"

	"github.com/smgladkovskiy/file_storage/internal/pkg/file"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCombine(t *testing.T) {
	type expStruct struct {
		data []byte
		err  error
	}

	type testCase struct {
		name string
		in   [][]byte
		exp  expStruct
	}

	testCases := []testCase{
		{
			name: "ordinary combine",
			in:   [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			exp: expStruct{
				data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
				err:  nil,
			},
		},
		{
			name: "combine with nil tail",
			in:   [][]byte{{1, 2, 3}, {4, 5, 6}, nil},
			exp: expStruct{
				data: []byte{1, 2, 3, 4, 5, 6},
				err:  nil,
			},
		},
		{
			name: "combine different parts length",
			in:   [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8}},
			exp: expStruct{
				data: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				err:  nil,
			},
		},
		{
			name: "nil parts",
			in:   nil,
			exp: expStruct{
				data: nil,
				err:  file.ErrEmptyParts,
			},
		},
		{
			name: "nil parts 2",
			in:   [][]byte{nil, nil, nil},
			exp: expStruct{
				data: nil,
				err:  nil,
			},
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out, err := file.Combine(tc.in)
			if tc.exp.err != nil {
				require.ErrorIs(t, tc.exp.err, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp.data, out)
		})
	}
}
