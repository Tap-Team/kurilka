package filestorage

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/google/uuid"
)

const _PROVIDER = "sharestatic/filestorage.fileStorage"

type FileStorage interface {
	SaveFile(ctx context.Context, data []byte) (string, error)
}

type fileStorage struct {
	storagePath       string
	fileExpiration    time.Duration
	expirationStorage FileExpirationStorage
}

func NewFileStorage(ctx context.Context, filePath string, fileExpiration time.Duration) FileStorage {
	os.Mkdir(filePath, os.ModePerm)
	storage := &fileStorage{storagePath: filePath, fileExpiration: fileExpiration, expirationStorage: NewFileExpirationStorage()}
	go storage.run(ctx)
	return storage
}

func (f *fileStorage) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			now := time.Now()
			go f.clear(ctx, now)
			time.Sleep(time.Second)
		}
	}
}

func (f *fileStorage) clear(ctx context.Context, now time.Time) {
	files := f.expirationStorage.ClearFiles(ctx, now)
	for _, name := range files {
		filePath := path.Join(f.storagePath, name)
		os.Remove(filePath)
	}
}

func (f *fileStorage) SaveFile(ctx context.Context, data []byte) (string, error) {
	name := uuid.New().String() + ".png"
	file, err := os.Create(path.Join(f.storagePath, name))
	if err != nil {
		return "", exception.Wrap(err, exception.NewCause("failed create file", "SaveFile", _PROVIDER))
	}
	_, err = file.Write(data)
	if err != nil {
		return "", exception.Wrap(err, exception.NewCause("failed write data to file", "SaveFile", _PROVIDER))
	}
	f.expirationStorage.AddFiles(ctx, time.Now().Add(f.fileExpiration), name)
	return name, nil
}
