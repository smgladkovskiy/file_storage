package file_test

import (
	"testing"

	"github.com/smgladkovskiy/file_storage/internal/pkg/file"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	type inStruct struct {
		slice []byte
		parts int
	}
	type testCase struct {
		name string
		in   inStruct
		exp  [][]byte
	}

	testCases := []testCase{
		{
			name: "10 bytes to 3 parts",
			in: inStruct{
				slice: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				parts: 3,
			},
			exp: [][]byte{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}},
		},
		{
			name: "3 bytes to 3 parts",
			in: inStruct{
				slice: []byte{1, 2, 3},
				parts: 3,
			},
			exp: [][]byte{{1}, {2}, {3}},
		},
		{
			name: "3 bytes to 5 parts",
			in: inStruct{
				slice: []byte{1, 2, 3},
				parts: 5,
			},
			exp: [][]byte{{1}, {2}, {3}, nil, nil},
		},
		{
			name: "0 bytes to 5 parts",
			in: inStruct{
				slice: []byte{},
				parts: 5,
			},
			exp: nil,
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := file.Split(tc.in.slice, tc.in.parts)
			assert.Equal(t, tc.exp, out)
		})
	}
}
