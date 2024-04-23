package main

import (
	"github.com/AustrianDataLAB/GeWoScout/demo/config"
	"github.com/AustrianDataLAB/GeWoScout/demo/controller"
	"github.com/AustrianDataLAB/GeWoScout/demo/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Create HTTP server
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"hello": "world",
		})
	})

	// Connect To Database
	config.DatabaseInit()
	gorm := config.DB()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()

	// Example of using the database object
	if err := gorm.AutoMigrate(&model.Book{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	bookRoute := e.Group("/book")
	bookRoute.POST("/", controller.CreateBook)
	bookRoute.GET("/:id", controller.GetBook)
	bookRoute.PUT("/:id", controller.UpdateBook)
	bookRoute.DELETE("/:id", controller.DeleteBook)

	e.Logger.Fatal(e.Start(":8080"))
}
