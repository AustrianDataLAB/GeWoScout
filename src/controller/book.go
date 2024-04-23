package controller

import (
	"errors"
	"github.com/AustrianDataLAB/GeWoScout/demo/config"
	"github.com/AustrianDataLAB/GeWoScout/demo/model"
	"gorm.io/gorm"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateBook(c echo.Context) error {
	b := new(model.Book)
	db := config.DB()

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	if err := db.Create(&b).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": b,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateBook(c echo.Context) error {
	id := c.Param("id")
	b := new(model.Book)
	db := config.DB()

	// Binding data
	if err := c.Bind(b); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	existing_book := new(model.Book)

	if err := db.First(&existing_book, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	existing_book.Name = b.Name
	existing_book.Description = b.Description
	if err := db.Save(&existing_book).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_book,
	}

	return c.JSON(http.StatusOK, response)
}

func GetBook(c echo.Context) error {
	id := c.Param("id")
	db := config.DB()

	var books []*model.Book

	if res := db.First(&books, id); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			http.NotFound(c.Response(), c.Request())
			return nil
		}
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusOK, data)
	}

	response := map[string]interface{}{
		"data": books[0],
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteBook(c echo.Context) error {
	id := c.Param("id")
	db := config.DB()

	book := new(model.Book)

	err := db.Delete(&book, id).Error
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "a book has been deleted",
	}
	return c.JSON(http.StatusOK, response)
}
