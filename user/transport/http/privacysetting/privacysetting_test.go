package privacysetting_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/privacysettingerror"
	"github.com/Tap-Team/kurilka/internal/errorutils/userprivacysettingerror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/privacysettingdatamanager"
	"github.com/Tap-Team/kurilka/user/transport/http/privacysetting"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestGetPrivacySetttingsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	manager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)

	transport := privacysetting.NewPrivacySettingTransport(manager)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		privacySettings []usermodel.PrivacySetting
		err             error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},

		{
			userId: 12,
			queryValues: map[string]string{
				"vk_user_id": "12",
			},
			err:        userprivacysettingerror.ExceptionUserPrivacySettingNotFound(),
			statusCode: http.StatusNotFound,
		},
		{
			userId: 1,
			queryValues: map[string]string{
				"vk_user_id": "1",
			},
			privacySettings: []usermodel.PrivacySetting{
				usermodel.ACHIEVEMENTS_CIGARETTE,
				usermodel.ACHIEVEMENTS_SAVING,
				usermodel.ACHIEVEMENTS_DURATION,
			},
			statusCode: http.StatusOK,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			manager.EXPECT().PrivacySettings(gomock.Any(), cs.userId).Return(cs.privacySettings, cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/privacysettings?"+urlValues.Encode(), nil)

		transport.GetPrivacySettingsHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		} else {
			privacySettings := make([]usermodel.PrivacySetting, 0, len(cs.privacySettings))
			err := json.NewDecoder(rec.Result().Body).Decode(&privacySettings)
			rec.Result().Body.Close()
			assert.NilError(t, err)

			equal := slices.Equal(privacySettings, cs.privacySettings)
			assert.Equal(t, true, equal, "privacy settings not equal")
		}
	}
}

func TestRemovePrivacySettingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	manager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)

	transport := privacysetting.NewPrivacySettingTransport(manager)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		privacySetting usermodel.PrivacySetting
		err            error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},

		{
			queryValues: map[string]string{
				"vk_user_id":     "1",
				"privacySetting": "asdfasdfa",
			},
			err:        privacysettingerror.ExceptionPrivacySettingNotExist(),
			statusCode: http.StatusBadRequest,
		},
		{
			userId: 12,
			queryValues: map[string]string{
				"vk_user_id":     "12",
				"privacySetting": string(usermodel.STATISTICS_CIGARETTE),
			},
			privacySetting: usermodel.STATISTICS_CIGARETTE,
			err:            userprivacysettingerror.ExceptionUserPrivacySettingNotFound(),
			statusCode:     http.StatusNotFound,
		},

		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id":     "123",
				"privacySetting": string(usermodel.STATISTICS_MONEY),
			},
			privacySetting: usermodel.STATISTICS_MONEY,
			statusCode:     http.StatusNoContent,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			manager.EXPECT().Remove(gomock.Any(), cs.userId, cs.privacySetting).Return(cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/privacysettings/remove?"+urlValues.Encode(), nil)

		transport.RemovePrivacySettingHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}
}

func TestAddPrivacySettingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()

	manager := privacysettingdatamanager.NewMockPrivacySettingManager(ctrl)

	transport := privacysetting.NewPrivacySettingTransport(manager)

	cases := []struct {
		userId      int64
		queryValues map[string]string

		privacySetting usermodel.PrivacySetting
		err            error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			queryValues: map[string]string{
				"vk_user_id":     "1",
				"privacySetting": "asdfasdfa",
			},
			err:        privacysettingerror.ExceptionPrivacySettingNotExist(),
			statusCode: http.StatusBadRequest,
		},
		{
			userId: 12,
			queryValues: map[string]string{
				"vk_user_id":     "12",
				"privacySetting": string(usermodel.STATISTICS_CIGARETTE),
			},
			privacySetting: usermodel.STATISTICS_CIGARETTE,
			err:            userprivacysettingerror.ExceptionUserPrivacySettingNotFound(),
			statusCode:     http.StatusNotFound,
		},

		{
			userId: 123,
			queryValues: map[string]string{
				"vk_user_id":     "123",
				"privacySetting": string(usermodel.STATISTICS_MONEY),
			},
			privacySetting: usermodel.STATISTICS_MONEY,
			statusCode:     http.StatusNoContent,
		},
	}

	for _, cs := range cases {
		if cs.userId != 0 {
			manager.EXPECT().Add(gomock.Any(), cs.userId, cs.privacySetting).Return(cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/privacysettings/remove?"+urlValues.Encode(), nil)

		transport.AddRemovePrivacySettingHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}
}
