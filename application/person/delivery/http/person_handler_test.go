package http

import (
	"bytes"
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person/mocks"
	server "github.com/go-park-mail-ru/2020_1_k-on/application/server/mocks"
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

var testPerson = models.Person{
	Id:         1,
	Name:       "test",
	Occupation: "test",
	BirthDate:  "2020-07-22",
	BirthPlace: "test",
	Image:      "test",
	Films:      nil,
	Series:     nil,
}

const (
	id      = "1"
	wrongId = "abc"
	param   = "id"
	ok      = 200
)

func beforeTest(t *testing.T) (*PersonHandler, *mocks.MockUseCase, *server.MockContext, *server.MockResponseWriter) {
	ctrl := gomock.NewController(t)
	w := server.NewMockResponseWriter(ctrl)
	ctx := server.NewMockContext(ctrl)
	usecase := mocks.NewMockUseCase(ctrl)
	handler := PersonHandler{usecase: usecase, logger: zap.NewExample(), sanitizer: bluemonday.UGCPolicy()}

	response := echo.NewResponse(w, echo.New())
	ctx.EXPECT().Response().Return(response).AnyTimes()
	w.EXPECT().Write(gomock.Any()).AnyTimes()

	return &handler, usecase, ctx, w
}

func TestPersonHandler_GetById(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(id)
	usecase.EXPECT().GetById(testPerson.Id).Return(&testPerson, nil)
	w.EXPECT().WriteHeader(ok)

	err = handler.GetById(ctx)

	assert.Nil(t, err)
}

func TestPersonHandler_GetById_WrongParameter(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)
	request, err := http.NewRequest("", "", bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	ctx.EXPECT().Param(param).Return(wrongId)
	w.EXPECT().WriteHeader(http.StatusBadRequest)

	err = handler.GetById(ctx)

	assert.NotNil(t, err)
}

func TestPersonHandler_Add(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)

	body, err := easyjson.Marshal(testPerson)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	usecase.EXPECT().Add(&testPerson).Return(&testPerson, nil)
	w.EXPECT().WriteHeader(ok)

	err = handler.Add(ctx)

	assert.Nil(t, err)
}

func TestPersonHandler_Add_BadRequest(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)

	request, err := http.NewRequest("", "", strings.NewReader(""))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	w.EXPECT().WriteHeader(http.StatusBadRequest)

	err = handler.Add(ctx)

	assert.NotNil(t, err)
}

func TestPersonHandler_Update(t *testing.T) {
	handler, usecase, ctx, w := beforeTest(t)

	body, err := easyjson.Marshal(testPerson)
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}
	request, err := http.NewRequest("", "", bytes.NewReader(body))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	usecase.EXPECT().Update(&testPerson).Return(&testPerson, nil)
	w.EXPECT().WriteHeader(ok)

	err = handler.Update(ctx)

	assert.Nil(t, err)
}

func TestPersonHandler_Update_BadRequest(t *testing.T) {
	handler, _, ctx, w := beforeTest(t)

	request, err := http.NewRequest("", "", strings.NewReader(""))
	if err != nil {
		t.Errorf("unexpected error: '%s'", err)
	}

	ctx.EXPECT().Request().Return(request).AnyTimes()
	w.EXPECT().WriteHeader(http.StatusBadRequest)

	err = handler.Update(ctx)

	assert.NotNil(t, err)
}
