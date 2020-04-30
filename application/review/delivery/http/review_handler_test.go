package http

import (
	"bytes"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/review/mocks"
	server "github.com/go-park-mail-ru/2020_1_k-on/application/server/mocks"
	"github.com/go-park-mail-ru/2020_1_k-on/pkg/constants"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"testing"
)

var testReview = models.Review{
	Id:        1,
	Rating:    10,
	Body:      "test",
	UserId:    1,
	ProductId: 1,
	Usr:       models.ListUser{},
}

const (
	id      = "1"
	wrongId = "abc"
	param   = "id"
	ok      = 200
)

func beforeTest(t *testing.T) (*ReviewHandler, *mocks.MockUseCase, *server.MockContext, *server.MockResponseWriter) {
	ctrl := gomock.NewController(t)
	w := server.NewMockResponseWriter(ctrl)
	ctx := server.NewMockContext(ctrl)
	usecase := mocks.NewMockUseCase(ctrl)
	handler := ReviewHandler{series: usecase, film: usecase, logger: zap.NewExample(), sanitizer: bluemonday.UGCPolicy()}

	response := echo.NewResponse(w, echo.New())
	ctx.EXPECT().Response().Return(response).AnyTimes()

	return &handler, usecase, ctx, w
}

func TestReviewHandler_AddFilmReview(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)

	body, err := easyjson.Marshal(testReview)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Get(constants.UserIdKey).Return(testReview.UserId)
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().Add(gomock.Any()).Return(nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = handler.AddFilmReview(ctx)
	assert.Nil(t, err)
}

func TestReviewHandler_AddFilmReview_BadRequest(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)

	request, err := http.NewRequest("", "", strings.NewReader(""))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(id)
	w.EXPECT().WriteHeader(http.StatusBadRequest)
	w.EXPECT().Write(gomock.Any())

	err = handler.AddFilmReview(ctx)
	assert.NotNil(t, err)
}

func TestReviewHandler_AddSeriesReview(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)

	body, err := easyjson.Marshal(testReview)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Get(constants.UserIdKey).Return(testReview.UserId)
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().Add(gomock.Any()).Return(nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = handler.AddSeriesReview(ctx)
	assert.Nil(t, err)
}

func TestReviewHandler_AddSeriesReview_BadRequest(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)

	request, err := http.NewRequest("", "", strings.NewReader(""))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(id)
	w.EXPECT().WriteHeader(http.StatusBadRequest)
	w.EXPECT().Write(gomock.Any())

	err = handler.AddSeriesReview(ctx)
	assert.NotNil(t, err)
}

func TestReviewHandler_GetByFilm(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().GetByProductId(testReview.ProductId).Return([]models.Review{testReview}, nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetByFilm(ctx)
	assert.Nil(t, err)
}

func TestReviewHandler_GetByFilm_WrongParameter(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(wrongId)
	w.EXPECT().WriteHeader(http.StatusBadRequest)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetByFilm(ctx)
	assert.NotNil(t, err)
}

func TestReviewHandler_GetBySeries(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().GetByProductId(testReview.ProductId).Return([]models.Review{testReview}, nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetBySeries(ctx)
	assert.Nil(t, err)
}

func TestReviewHandler_GetBySeries_WrongParameter(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(wrongId)
	w.EXPECT().WriteHeader(http.StatusBadRequest)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetBySeries(ctx)
	assert.NotNil(t, err)
}

func TestReviewHandler_GetByFilmAndUser(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Get(constants.UserIdKey).Return(testReview.UserId)
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().GetReview(testReview.ProductId, testReview.UserId).Return(&testReview, nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetByFilmAndUser(ctx)
	assert.Nil(t, err)
}

func TestReviewHandler_GetByFilmAndUser_WrongParameter(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Get(constants.UserIdKey).Return(testReview.UserId)
	ctx.EXPECT().Param(param).Return(wrongId)
	w.EXPECT().WriteHeader(http.StatusBadRequest)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetByFilmAndUser(ctx)
	assert.NotNil(t, err)
}

func TestReviewHandler_GetBySeriesAndUser(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Get(constants.UserIdKey).Return(testReview.UserId)
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().GetReview(testReview.ProductId, testReview.UserId).Return(&testReview, nil)
	w.EXPECT().WriteHeader(ok)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetBySeriesAndUser(ctx)
	assert.Nil(t, err)
}

func TestReviewHandler_GetBySeriesAndUser_WrongParameter(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Get(constants.UserIdKey).Return(testReview.UserId)
	ctx.EXPECT().Param(param).Return(wrongId)
	w.EXPECT().WriteHeader(http.StatusBadRequest)
	w.EXPECT().Write(gomock.Any())

	err = handler.GetByFilmAndUser(ctx)
	assert.NotNil(t, err)
}
