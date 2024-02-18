package controller_test

import (
	"dcard-backend/controller"
	"dcard-backend/domain"
	"dcard-backend/domain/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostAd_WrongJSONProvided_ShouldReturnBadRequestError(t *testing.T) {
	mockAdUsecase := mocks.NewAdUsecase(t)

	testAdController := controller.AdController{
		AdUsecase: mockAdUsecase,
	}

	reader := strings.NewReader(`{"wrong_format: value}`)

	httpRecorder := httptest.NewRecorder()
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/ad", reader)
	httpRequest.Header.Set("Content-Type", "application/json")

	app := gin.Default()
	app.POST("/api/v1/ad", testAdController.PostAd)
	app.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
}

func TestPostAd_CreateFail_ShouldReturnInternalServerError(t *testing.T) {
	mockAd := domain.Ad{
		Title: "Test AD",
	}
	mockAdUsecase := mocks.NewAdUsecase(t)
	mockAdUsecase.On("Create", mock.Anything, &mockAd).Return(errors.New("Fail")).Once()

	testAdController := controller.AdController{
		AdUsecase: mockAdUsecase,
	}

	reader := strings.NewReader(`{"title": "Test AD"}`)

	httpRecorder := httptest.NewRecorder()
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/ad", reader)
	httpRequest.Header.Set("Content-Type", "application/json")

	app := gin.Default()
	app.POST("/api/v1/ad", testAdController.PostAd)
	app.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusInternalServerError, httpRecorder.Code)
}

func TestPostAd_CreateSucess_ShouldReturnOK(t *testing.T) {
	mockAd := domain.Ad{
		Title: "Test AD",
	}
	mockAdUsecase := mocks.NewAdUsecase(t)
	mockAdUsecase.On("Create", mock.Anything, &mockAd).Return(nil).Once()

	testAdController := controller.AdController{
		AdUsecase: mockAdUsecase,
	}

	reader := strings.NewReader(`{"title": "Test AD"}`)

	httpRecorder := httptest.NewRecorder()
	httpRequest := httptest.NewRequest(http.MethodPost, "/api/v1/ad", reader)
	httpRequest.Header.Set("Content-Type", "application/json")

	app := gin.Default()
	app.POST("/api/v1/ad", testAdController.PostAd)
	app.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusOK, httpRecorder.Code)
}

func TestGetAdWithCondition_Fail_ShouldReturnInternalServerError(t *testing.T) {
	mockCondition := map[string][]string{}

	mockAdUsecase := mocks.NewAdUsecase(t)
	mockAdUsecase.On("GetByCondition", mock.Anything, mockCondition).Return([]domain.Ad{}, errors.New("Fail")).Once()

	testAdController := controller.AdController{
		AdUsecase: mockAdUsecase,
	}

	httpRecorder := httptest.NewRecorder()
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/ad", nil)

	app := gin.Default()
	app.GET("/api/v1/ad", testAdController.GetAdWithCondition)
	app.ServeHTTP(httpRecorder, httpRequest)

	assert.Equal(t, http.StatusInternalServerError, httpRecorder.Code)
}

func TestGetAdWithCondition_Success_ShouldReturnAds(t *testing.T) {
	mockCondition := map[string][]string{"offset": {"0"}}
	mockAds := []domain.Ad{
		{
			Title: "TEST AD",
		},
	}

	mockAdUsecase := mocks.NewAdUsecase(t)
	mockAdUsecase.On("GetByCondition", mock.Anything, mockCondition).Return(mockAds, nil).Once()

	testAdController := controller.AdController{
		AdUsecase: mockAdUsecase,
	}

	httpRecorder := httptest.NewRecorder()
	httpRequest := httptest.NewRequest(http.MethodGet, "/api/v1/ad?offset=0", nil)

	app := gin.Default()
	app.GET("/api/v1/ad", testAdController.GetAdWithCondition)
	app.ServeHTTP(httpRecorder, httpRequest)

	var responseAds map[string][]domain.Ad
	err := json.Unmarshal(httpRecorder.Body.Bytes(), &responseAds)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, httpRecorder.Code)
	assert.EqualValues(t, mockAds, responseAds["items"])
}
