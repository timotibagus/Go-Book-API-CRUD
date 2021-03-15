package main

import (
	"book-api/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	book := r.Group("/books")
	{
		book.GET("/", controller.BookGet)
		book.GET("/:id", controller.BookByIdGet)
		book.POST("/", controller.BookPost)
		book.DELETE("/:id", controller.BookDelete)
		book.PATCH("/:id", controller.BookUpate)
	}

	r.Run()
}