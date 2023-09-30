package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/sharestatic/filestorage"
)

const _PROVIDER = "sharestatic/handler.Handler"

type Handler struct {
	fileStorage filestorage.FileStorage
	urlPrefix   string
}

func New(fileStorage filestorage.FileStorage, urlPrefix string) *Handler {
	return &Handler{fileStorage: fileStorage, urlPrefix: urlPrefix}
}

func (h *Handler) PushImageHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			httphelpers.Error(w, err)
			return
		}
		fileName, err := h.fileStorage.SaveFile(ctx, b)
		if err != nil {
			err := exception.Wrap(err, exception.NewCause("save file in file storage", "PushImageHandler", _PROVIDER))
			httphelpers.Error(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(h.urlPrefix + fileName))
	}
	return http.HandlerFunc(handler)
}
