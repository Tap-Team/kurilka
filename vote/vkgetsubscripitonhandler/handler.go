package vkgetsubscripitonhandler

import (
	"io"
	"net/http"
	"strings"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
)

type VKGetSubscriptionHandler struct {
	client        *http.Client
	Version       string
	AppServiceKey string
}

func New(client *http.Client, version, appServiceKey string) http.Handler {
	return &VKGetSubscriptionHandler{
		client:        client,
		Version:       version,
		AppServiceKey: appServiceKey,
	}
}

func (v *VKGetSubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	values := r.URL.Query()
	values.Set("v", v.Version)
	values.Set("access_token", v.AppServiceKey)
	values.Set("user_id", values.Get("vk_user_id"))
	request, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.vk.com/method/orders.getUserSubscriptionById", strings.NewReader(values.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := v.client.Do(request)
	if err != nil {
		httphelpers.Error(w, err)
		return
	}
	io.Copy(w, resp.Body)
	resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
}
