package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dcard-backend/domain"
)

type AdController struct {
	AdUsecase domain.AdUsecase
}

// PostAd       godoc
// @Summary     Admin API
// @Description Create an ad
// @Tags        ad
// @Accept      json
// @Produce     json
// @Param       ad body domain.Ad True "Add an ad"
// @Success     200 {object} domain.SuccessResponse
// @Failure     400 {object} domain.ErrorResponse
// @Failure     500 {object} domain.ErrorResponse
// @Router      /ad [post]
func (ac *AdController) PostAd(ctx *gin.Context) {
	var ad domain.Ad
	if err := ctx.Bind(&ad); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := ac.AdUsecase.Create(ctx.Request.Context(), &ad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, domain.SuccessResponse{Message: "Ad insert successfully"})
}

// GetAdWithCondition godoc
// @Summary           Public API
// @Description       Get a list of ads with queries
// @Tags              ad
// @Produce           json
// @Param             offset   query int    true  "Get ads starting from offset"
// @Param             limit    query int    false "Get how many ads" default(5)
// @Param             age      query int    false "Target age"
// @Param             gender   query int    false "Target gender"
// @Param             country  query string false "Target country"
// @Param             platform query string false "Target platform"
// @Success           200 {object} map[string][]domain.Ad "{\"items\": [ad, ...]}"
// @Failure           500 {object} domain.ErrorResponse
// @Router            /ad [get]
func (ac *AdController) GetAdWithCondition(ctx *gin.Context) {
	condition := ctx.Request.URL.Query()

	ads, err := ac.AdUsecase.GetByCondition(ctx.Request.Context(), condition)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": ads,
	})
}
