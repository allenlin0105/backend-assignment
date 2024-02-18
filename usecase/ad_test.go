package usecase_test

import (
	"context"
	"dcard-backend/domain"
	"dcard-backend/domain/mocks"
	"dcard-backend/usecase"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_TitleNotProvided_ShouldReturnError(t *testing.T) {
	mockAd := domain.Ad{
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2024-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)
	assert.Error(t, err)
}

func TestCreate_StartAtNotProvided_ShouldReturnError(t *testing.T) {
	mockAd := domain.Ad{
		Title: "Test AD",
		EndAt: "2024-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)
	assert.Error(t, err)
}

func TestCreate_EndAtNotProvided_ShouldReturnError(t *testing.T) {
	mockAd := domain.Ad{
		Title: "Test AD",
		EndAt: "2024-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)
	assert.Error(t, err)
}

func TestCreate_ConditionNotProvided_ShouldCreateEmptyCondition(t *testing.T) {
	mockAd := domain.Ad{
		Title:   "Test AD",
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2025-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("Create", mock.Anything, &mockAd).Return(nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)

	assert.NoError(t, err)
	assert.NotNil(t, mockAd.Condition)
}

func TestCreate_StartAtEndAtProvided_ShouldChangeToUTF8(t *testing.T) {
	mockAd := domain.Ad{
		Title:   "Test AD",
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2025-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("Create", mock.Anything, &mockAd).Return(nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	startAtUTF0, _ := time.Parse(time.RFC3339, mockAd.StartAt)
	endAtUTF0, _ := time.Parse(time.RFC3339, mockAd.EndAt)
	_, startAtUTF0Offset := startAtUTF0.Zone()
	_, endAtUTF0offset := endAtUTF0.Zone()

	err := testAdUsecase.Create(context.Background(), &mockAd)

	startAtUTF8, _ := time.Parse(time.RFC3339, mockAd.StartAt)
	endAtUTF8, _ := time.Parse(time.RFC3339, mockAd.EndAt)
	_, startAtUTF8Offset := startAtUTF8.Zone()
	_, endAtUTF8Offset := endAtUTF8.Zone()

	assert.NoError(t, err)
	assert.EqualValues(t, startAtUTF0, startAtUTF8.UTC())
	assert.EqualValues(t, endAtUTF0, endAtUTF8.UTC())
	assert.Equal(t, 8*60*60, startAtUTF8Offset-startAtUTF0Offset)
	assert.Equal(t, 8*60*60, endAtUTF8Offset-endAtUTF0offset)
}

func TestCreate_AgeStartNotProvided_ShouldSetTo1(t *testing.T) {
	mockAd := domain.Ad{
		Title:   "Test AD",
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2025-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("Create", mock.Anything, &mockAd).Return(nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)

	assert.NoError(t, err)
	assert.Equal(t, 1, mockAd.Condition.AgeStart)
}

func TestCreate_AgeEndNotProvided_ShouldSetTo100(t *testing.T) {
	mockAd := domain.Ad{
		Title:   "Test AD",
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2025-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("Create", mock.Anything, &mockAd).Return(nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)

	assert.NoError(t, err)
	assert.Equal(t, 100, mockAd.Condition.AgeEnd)
}

func TestCreate_GenderAndCountryAndPlatformNotProvided_ShouldAddAnyValueToSlice(t *testing.T) {
	mockAd := domain.Ad{
		Title:   "Test AD",
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2025-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("Create", mock.Anything, &mockAd).Return(nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)

	assert.NoError(t, err)
	assert.Equal(t, []string{"A"}, mockAd.Condition.Gender)
	assert.Equal(t, []string{"AY"}, mockAd.Condition.Country)
	assert.Equal(t, []string{"any"}, mockAd.Condition.Platform)
}

func TestCreate_AdRepositoryFail_ShouldReturnError(t *testing.T) {
	mockAd := domain.Ad{
		Title:   "Test AD",
		StartAt: "2024-01-01T00:00:00.000Z",
		EndAt:   "2025-01-01T00:00:00.000Z",
	}
	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("Create", mock.Anything, &mockAd).Return(errors.New("Fail")).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	err := testAdUsecase.Create(context.Background(), &mockAd)

	assert.Error(t, err)
}

func TestGetByCondition_IfOffsetNotProvided_ShouldReturnError(t *testing.T) {
	condition := map[string][]string{}

	mockAdRepository := mocks.NewAdRepository(t)

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	_, err := testAdUsecase.GetByCondition(context.Background(), condition)

	assert.Error(t, err)
}

func TestGetByCondition_IfLimitNotProvided_ShouldAssignValue5(t *testing.T) {
	condition := map[string][]string{
		"offset": {"0"},
	}

	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("GetByCondition", mock.Anything, condition).Return([]domain.Ad{}, nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	testAdUsecase.GetByCondition(context.Background(), condition)

	if assert.NotNil(t, condition["limit"]) {
		assert.Equal(t, []string{"5"}, condition["limit"])
	}
}

func TestGetByCondition_IfConditionProvided_ShouldAppendAnyValueToSlice(t *testing.T) {
	condition := map[string][]string{
		"offset":   {"0"},
		"gender":   {"M"},
		"country":  {"TW"},
		"platform": {"web"},
	}

	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("GetByCondition", mock.Anything, condition).Return([]domain.Ad{}, nil).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	testAdUsecase.GetByCondition(context.Background(), condition)

	assert.Equal(t, []string{"M", "A"}, condition["gender"])
	assert.Equal(t, []string{"TW", "AY"}, condition["country"])
	assert.Equal(t, []string{"web", "any"}, condition["platform"])
}

func TestGetByCondition_IfAdRepositoryFail_ShouldReturnError(t *testing.T) {
	condition := map[string][]string{
		"offset": {"0"},
	}

	mockAdRepository := mocks.NewAdRepository(t)
	mockAdRepository.On("GetByCondition", mock.Anything, condition).Return([]domain.Ad{}, errors.New("Fail")).Once()

	testAdUsecase := usecase.NewAdUsecase(mockAdRepository, time.Second*1)

	_, err := testAdUsecase.GetByCondition(context.Background(), condition)

	assert.Error(t, err)
}
