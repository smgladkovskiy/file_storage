package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/smgladkovskiy/file_storage/internal/domain/entity"

	"github.com/google/uuid"
)

type Storage struct {
	files          map[uuid.UUID]entity.FileMeta
	storagesAmount int
	storagesMap    map[int]struct{}
	storages       []map[uuid.UUID][]byte

	filesMu       sync.RWMutex
	storagesMapMu sync.RWMutex
	storagesMu    sync.RWMutex
}

var (
	ErrNotFound                   = errors.New("file not found")
	ErrFileChunksMoreThanStorages = errors.New("file chunks are more than file storages")
)

func NewInMemoryStorage(storagesAmount int) *Storage {
	s := &Storage{
		storagesAmount: storagesAmount,
		files:          make(map[uuid.UUID]entity.FileMeta),
		storagesMap:    make(map[int]struct{}),
		storages:       make([]map[uuid.UUID][]byte, storagesAmount),

		filesMu:       sync.RWMutex{},
		storagesMapMu: sync.RWMutex{},
		storagesMu:    sync.RWMutex{},
	}
	for i := 0; i < storagesAmount; i++ {
		s.storages[i] = make(map[uuid.UUID][]byte)
		s.storagesMap[i] = struct{}{}
	}

	return s
}

func (s *Storage) SaveFile(_ context.Context, file *entity.File) error {
	if err := file.Split(s.storagesAmount); err != nil {
		return err
	}

	s.filesMu.RLock()
	s.files[file.ID] = file.FileMeta
	s.filesMu.RUnlock()

	if len(file.SplitContent) > s.storagesAmount {
		return ErrFileChunksMoreThanStorages
	}

	s.storagesMapMu.RLock()
	defer s.storagesMapMu.RUnlock()

	wg := sync.WaitGroup{}
	for i := range s.storagesMap {
		if i > len(file.SplitContent)-1 {
			continue
		}

		wg.Add(1)
		go func(id int) {
			s.storeFileChunk(id, file.ID, file.SplitContent[id])
			wg.Done()
		}(i)
	}
	wg.Wait()

	return nil
}

func (s *Storage) storeFileChunk(storageID int, fileID uuid.UUID, content []byte) {
	s.storagesMu.RLock()
	_, ok := s.storages[storageID][fileID]
	s.storagesMu.RUnlock()

	s.storagesMu.Lock()
	defer s.storagesMu.Unlock()

	if !ok {
		s.storages[storageID][fileID] = make([]byte, len(content))
	}

	s.storages[storageID][fileID] = content
}

func (s *Storage) GetFileByID(_ context.Context, fileID uuid.UUID) (*entity.File, error) {
	s.filesMu.RLock()
	fileMeta, ok := s.files[fileID]
	s.filesMu.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	file := entity.NewFileFromMeta(fileID, fileMeta)

	parts := make([][]byte, fileMeta.Parts)
	s.storagesMu.RLock()
	for i := range s.storages {
		part, ok := s.storages[i][fileID]
		if !ok {
			continue
		}

		parts = append(parts, part)
	}
	s.storagesMu.RUnlock()

	if err := file.SetSplitContent(parts); err != nil {
		return nil, err
	}

	return file, nil
}
