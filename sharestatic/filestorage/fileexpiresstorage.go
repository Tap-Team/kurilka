package filestorage

import (
	"context"
	"sync"
	"time"
)

type FileExpirationStorage interface {
	AddFiles(ctx context.Context, expired time.Time, files ...string)
	ClearFiles(ctx context.Context, expired time.Time) []string
}

type fileExpirationStorage struct {
	storage map[int64][]string
	mu      sync.Mutex
}

func NewFileExpirationStorage() FileExpirationStorage {
	return &fileExpirationStorage{
		storage: make(map[int64][]string),
	}
}

func (f *fileExpirationStorage) AddFiles(ctx context.Context, expired time.Time, files ...string) {
	f.mu.Lock()
	f.storage[expired.Unix()] = append(f.storage[expired.Unix()], files...)
	f.mu.Unlock()
}

func (f *fileExpirationStorage) ClearFiles(ctx context.Context, expired time.Time) []string {
	f.mu.Lock()
	files := f.storage[expired.Unix()]
	delete(f.storage, expired.Unix())
	f.mu.Unlock()
	return files
}
