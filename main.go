package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"dcard-backend/config"
	_ "dcard-backend/docs"
	"dcard-backend/router"
)

// @title  Dcard AD API
// @version 1.0
// @description The server for AD services

// @host 127.0.0.1:3000
// @BasePath /api/v1
func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.OpenMySQLDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer config.CloseMySQLDatabase(db)

	t, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	timeout := time.Duration(t) * time.Second

	app := gin.Default()
	app.Use(cors.Default())

	router.SetUpRoutes(app, db, timeout)

	port := os.Getenv("APP_PORT")
	app.Run(":" + port)
}
