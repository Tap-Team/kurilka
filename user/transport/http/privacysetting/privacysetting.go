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

const _PROVIDER = "user/transport/http/privacysetting"

type PrivacySettingTransport struct {
	privacySetting privacysettingdatamanager.PrivacySettingManager
}

func NewPrivacySettingTransport(privacySetting privacysettingdatamanager.PrivacySettingManager) *PrivacySettingTransport {
	return &PrivacySettingTransport{
		privacySetting: privacySetting,
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
//	@Param			vk_user_id	query	int64	true	"vk user id"
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
//	@Param			vk_user_id		query	int64						true	"vk user id"
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
//	@Param			vk_user_id		query	int64						true	"vk user id"
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
