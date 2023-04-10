package entity

import (
	"io"
	"io/fs"
)

type FileInterface interface {
	fs.FileInfo
	io.Reader
}
