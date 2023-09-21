package trigger_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Tap-Team/kurilka/internal/errorutils/triggererror"
	"github.com/Tap-Team/kurilka/internal/errorutils/usererror"
	"github.com/Tap-Team/kurilka/internal/httphelpers"
	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/user/datamanager/triggerdatamanager"
	"github.com/Tap-Team/kurilka/user/transport/http/trigger"
	"github.com/golang/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestTriggerRemoveHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := triggerdatamanager.NewMockTriggerManager(ctrl)
	transport := trigger.NewTriggerTransport(manager)

	cases := []struct {
		userId      int64
		trigger     usermodel.Trigger
		queryValues map[string]string

		err error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			queryValues: map[string]string{
				"vk_user_id": "1",
				"trigger":    "",
			},
			err:        triggererror.ExceptionTriggerNotExist(),
			statusCode: http.StatusBadRequest,
		},
		{
			userId:  123,
			trigger: usermodel.THANK_YOU,
			queryValues: map[string]string{
				"vk_user_id": "123",
				"trigger":    string(usermodel.THANK_YOU),
			},
			statusCode: http.StatusNoContent,
		},
	}

	for _, cs := range cases {

		if cs.userId != 0 {
			manager.EXPECT().Remove(gomock.Any(), cs.userId, cs.trigger).Return(cs.err).Times(1)
		}

		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/triggers/remove?"+urlValues.Encode(), nil)

		transport.RemoveTriggerHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}
}

func TestTriggerAddHandler(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	manager := triggerdatamanager.NewMockTriggerManager(ctrl)
	transport := trigger.NewTriggerTransport(manager)

	cases := []struct {
		userId      int64
		trigger     usermodel.Trigger
		queryValues map[string]string

		managerCall bool
		managerErr  error

		err error

		statusCode int
	}{
		{
			err:        httphelpers.ErrFailedParseVK_ID,
			statusCode: http.StatusInternalServerError,
		},
		{
			queryValues: map[string]string{
				"vk_user_id": "1",
				"trigger":    "",
			},
			err:        triggererror.ExceptionTriggerNotExist(),
			statusCode: http.StatusBadRequest,
		},
		{
			userId:  123,
			trigger: usermodel.THANK_YOU,
			queryValues: map[string]string{
				"vk_user_id": "123",
				"trigger":    string(usermodel.THANK_YOU),
			},
			managerCall: true,
			managerErr:  usererror.ExceptionUserNotFound(),
			err:         usererror.ExceptionUserNotFound(),

			statusCode: http.StatusNotFound,
		},
		{
			userId:  123,
			trigger: usermodel.THANK_YOU,
			queryValues: map[string]string{
				"vk_user_id": "123",
				"trigger":    string(usermodel.THANK_YOU),
			},
			managerCall: true,
			statusCode:  http.StatusNoContent,
		},
	}

	for _, cs := range cases {
		if cs.managerCall {
			manager.EXPECT().Add(gomock.Any(), cs.userId, cs.trigger).Return(cs.managerErr).Times(1)
		}
		urlValues := make(url.Values, len(cs.queryValues))
		for key, value := range cs.queryValues {
			urlValues.Set(key, value)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/triggers/remove?"+urlValues.Encode(), nil)

		transport.AddTriggerHandler(ctx).ServeHTTP(rec, req)

		assert.Equal(t, rec.Result().StatusCode, cs.statusCode, "wrong status code")

		if cs.err != nil {
			httphelpers.AssertError(t, cs.err, rec)
		}
	}
}
