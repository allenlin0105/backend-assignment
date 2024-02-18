package router

import (
	"database/sql"
	"dcard-backend/controller"
	"dcard-backend/repository"
	"dcard-backend/usecase"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(router *gin.Engine, db *sql.DB, timeout time.Duration) {
	ar := repository.NewAdRepository(db)
	au := usecase.NewAdUsecase(ar, timeout)
	ac := controller.AdController{
		AdUsecase: au,
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.POST("/api/v1/ad", ac.PostAd)
	router.GET("/api/v1/ad", ac.GetAdWithCondition)
}
