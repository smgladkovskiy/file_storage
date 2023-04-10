package entity

import (
	"errors"
	"fmt"
	"mime"
	"path/filepath"

	"github.com/smgladkovskiy/file_storage/internal/pkg/file"

	"github.com/google/uuid"
)

type FileMeta struct {
	Name  string
	Parts int
	Mime  string
}

type File struct {
	ID           uuid.UUID
	Content      []byte
	SplitContent [][]byte

	FileMeta
}

const minFileParts = 5

var (
	ErrEmptyFileName    = errors.New("empty file name")
	ErrNilFileID        = errors.New("nil file ID")
	ErrWrongPartsAmount = errors.New("wrong parts amount")
)

func NewFile(name string) (*File, error) {
	if name == "" {
		return nil, ErrEmptyFileName
	}

	return &File{
		ID: uuid.New(),
		FileMeta: FileMeta{
			Name: filepath.Base(name),
			Mime: mime.TypeByExtension(filepath.Ext(filepath.Base(name))),
		},
	}, nil
}

func NewFileFromID(id uuid.UUID) (*File, error) {
	if id == uuid.Nil {
		return nil, ErrNilFileID
	}

	return &File{ID: id}, nil
}

func NewFileFromMeta(id uuid.UUID, meta FileMeta) *File {
	return &File{ID: id, FileMeta: meta}
}

func (f *File) SetContent(content []byte) {
	f.Content = content
}

func (f *File) SetSplitContent(splitContent [][]byte) error {
	f.SplitContent = splitContent
	f.Parts = len(splitContent)

	return f.Combine()
}

func (f *File) Split(parts int) error {
	if parts < minFileParts {
		return fmt.Errorf("%w: got %d needs %d", ErrWrongPartsAmount, parts, minFileParts)
	}

	f.Parts = parts
	f.SplitContent = file.Split(f.Content, parts)

	return nil
}

func (f *File) Combine() error {
	var err error

	f.Content, err = file.Combine(f.SplitContent)

	return err
}
