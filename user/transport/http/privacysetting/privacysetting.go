package privacysetting

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/exception"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
)

//go:generate mockgen -source privacysetting.go -destination mocks.go -package privacysetting

const _PROVIDER = "user/transport/http/privacysetting"

type PrivacySettingSwitcher interface {
	Switch(ctx context.Context, userId int64, setting usermodel.PrivacySetting) error
}

type PrivacySettingTransport struct {
	privacySetting privacysettingdatamanager.PrivacySettingManager
	switcher       PrivacySettingSwitcher
}

func NewPrivacySettingTransport(privacySetting privacysettingdatamanager.PrivacySettingManager, switcher PrivacySettingSwitcher) *PrivacySettingTransport {
	return &PrivacySettingTransport{
		privacySetting: privacySetting,
		switcher:       switcher,
	}
}

type query url.Values

func (q query) PrivacySetting() usermodel.PrivacySetting {
	setting := url.Values(q).Get("privacySetting")
	return usermodel.PrivacySetting(setting)
}

// GetPrivacySettingHandler godoc
//
//	@Summary		GetPrivacySettingsHandler
//	@Description	get user privacy settings
//	@Tags			privacysettings
//	@Produce		json
//	@Success		200	{array}		usermodel.PrivacySetting
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/privacysettings [get]
func (t *PrivacySettingTransport) GetPrivacySettingsHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk user id", "GetPrivacySettingsHandler", _PROVIDER)))
			return
		}
		privacySettings, err := t.privacySetting.PrivacySettings(ctx, userId)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("get privacy settings", "GetPrivacySettingsHandler", _PROVIDER)))
			return
		}
		httphelpers.WriteJSON(w, privacySettings, http.StatusOK)
	}
	return http.HandlerFunc(handler)
}

// RemovePrivacySettingHandler
//
//	@Summary		RemovePrivacySetting
//	@Description	remove one user privacy setting, if setting not exists return error
//	@Tags			privacysettings
//	@Param			privacySetting	query	usermodel.PrivacySetting	true	"privacy setting"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/privacysettings/remove [delete]
func (t *PrivacySettingTransport) RemovePrivacySettingHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk user id", "RemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		privacySetting := query(r.URL.Query()).PrivacySetting()
		err = privacySetting.Validate()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("validate privacy setting", "RemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		err = t.privacySetting.Remove(ctx, userId, privacySetting)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("remove privacy setting", "RemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}

// AddPrivacySettingHandler
//
//	@Summary		AddPrivacySetting
//	@Description	add one user privacy setting, if setting exists return error
//	@Tags			privacysettings
//	@Param			privacySetting	query	usermodel.PrivacySetting	true	"privacy setting"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/privacysettings/add [post]
func (t *PrivacySettingTransport) AddRemovePrivacySettingHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk user id", "AddRemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		privacySetting := query(r.URL.Query()).PrivacySetting()
		err = privacySetting.Validate()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("validate privacy setting", "AddRemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		err = t.privacySetting.Add(ctx, userId, privacySetting)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("add privacy setting", "RemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	return http.HandlerFunc(handler)
}

// SwitchPrivacySettingHandler godoc
//
//	@Summary		SwitchPrivacySettings
//	@Description	add privacy setting if not exists and delete if exists
//	@Tags			privacysettings
//	@Param			privacySetting	query	usermodel.PrivacySetting	true	"privacy setting"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	errormodel.ErrorResponse
//	@Router			/privacysettings/switch [put]
func (h *PrivacySettingTransport) SwitchPrivacySettingHandler(ctx context.Context) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := httphelpers.VKID(r)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("parse vk user id", "AddRemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		privacySetting := query(r.URL.Query()).PrivacySetting()
		err = privacySetting.Validate()
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("validate privacy setting", "AddRemovePrivacySettingHandler", _PROVIDER)))
			return
		}
		err = h.switcher.Switch(ctx, userId, privacySetting)
		if err != nil {
			httphelpers.Error(w, exception.Wrap(err, exception.NewCause("switch privacy setting", "SwitchPrivacySettingHandler", _PROVIDER)))
			return
		}
		w.WriteHeader(http.StatusNoContent)

	}
	return http.HandlerFunc(handler)
}
