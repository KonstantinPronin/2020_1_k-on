package usecase

import (
	"github.com/go-park-mail-ru/2020_1_k-on/application/models"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person"
	"github.com/go-park-mail-ru/2020_1_k-on/application/person/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

var testListPerson = models.ListPerson{
	Id:   1,
	Name: "test",
}

func beforeTest(t *testing.T) (*mocks.MockRepository, person.UseCase) {
	ctrl := gomock.NewController(t)
	rep := mocks.NewMockRepository(ctrl)
	uc := NewPerson(rep, zap.NewExample())

	return rep, uc
}

func TestPerson_GetById(t *testing.T) {
	rep, uc := beforeTest(t)

	rep.EXPECT().GetById(testPerson.Id).Return(&testPerson, nil)

	actual, err := uc.GetById(testPerson.Id)

	assert.Nil(t, err)
	assert.Equal(t, testPerson, *actual)
}

func TestPerson_GetById_WrongId(t *testing.T) {
	_, uc := beforeTest(t)

	_, err := uc.GetById(0)

	assert.NotNil(t, err)
}

func TestPerson_Add(t *testing.T) {
	rep, uc := beforeTest(t)

	rep.EXPECT().Add(&testPerson).Return(&testPerson, nil)

	actual, err := uc.Add(&testPerson)

	assert.Nil(t, err)
	assert.Equal(t, testPerson, *actual)
}

func TestPerson_Add_EmptyName(t *testing.T) {
	_, uc := beforeTest(t)

	testPerson.Name = ""
	_, err := uc.Add(&testPerson)
	testPerson.Name = "test"

	assert.NotNil(t, err)
}

func TestPerson_Update(t *testing.T) {
	rep, uc := beforeTest(t)

	rep.EXPECT().Update(&testPerson).Return(&testPerson, nil)

	actual, err := uc.Update(&testPerson)

	assert.Nil(t, err)
	assert.Equal(t, testPerson, *actual)
}

func TestPerson_Update_WrongId(t *testing.T) {
	_, uc := beforeTest(t)

	testPerson.Id = 0
	_, err := uc.Update(&testPerson)
	testPerson.Id = 1

	assert.NotNil(t, err)
}

func TestPerson_GetActorsForFilm(t *testing.T) {
	rep, uc := beforeTest(t)
	id := uint(1)

	rep.EXPECT().GetActorsForFilm(id).Return(models.ListPersonArr{testListPerson}, nil)

	actual, err := uc.GetActorsForFilm(id)

	assert.Nil(t, err)
	assert.Equal(t, testListPerson, actual[0])
}

func TestPerson_GetActorsForFilm_WrongId(t *testing.T) {
	_, uc := beforeTest(t)
	id := uint(0)

	_, err := uc.GetActorsForFilm(id)

	assert.NotNil(t, err)
}

func TestPerson_GetActorsForSeries(t *testing.T) {
	rep, uc := beforeTest(t)
	id := uint(1)

	rep.EXPECT().GetActorsForSeries(id).Return(models.ListPersonArr{testListPerson}, nil)

	actual, err := uc.GetActorsForSeries(id)

	assert.Nil(t, err)
	assert.Equal(t, testListPerson, actual[0])
}

func TestPerson_GetActorsForSeries_WrongId(t *testing.T) {
	_, uc := beforeTest(t)
	id := uint(0)

	_, err := uc.GetActorsForSeries(id)

	assert.NotNil(t, err)
}
