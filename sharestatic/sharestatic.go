package sharestatic

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Tap-Team/kurilka/sharestatic/filestorage"
	"github.com/Tap-Team/kurilka/sharestatic/handler"
	"github.com/gorilla/mux"
)

type Config struct {
	Mux              *mux.Router
	ApiMux           *mux.Router
	FileSystemConfig struct {
		StoragePath    string
		FileExpiration time.Duration
	}
	StaticRouteConfig struct {
		StaticUrlPrefix string
		StaticRoute     string
	}
}

func disableListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func EnableShareStatic(cnf *Config) {
	ctx := context.Background()
	storagePath := cnf.FileSystemConfig.StoragePath
	urlPrefix := cnf.StaticRouteConfig.StaticUrlPrefix
	staticRoute := cnf.StaticRouteConfig.StaticRoute

	fileStorage := filestorage.NewFileStorage(ctx, storagePath, cnf.FileSystemConfig.FileExpiration)

	h := handler.New(fileStorage, urlPrefix)

	fileServer := http.FileServer(http.Dir(storagePath))
	r := cnf.Mux.NewRoute().Subrouter()
	r.Use(disableListing)
	r.PathPrefix(staticRoute).Handler(http.StripPrefix(staticRoute, fileServer))
	cnf.ApiMux.Handle("/push-share-static", h.PushImageHandler(ctx)).Methods(http.MethodPost)
}
