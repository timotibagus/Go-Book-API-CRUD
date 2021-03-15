package controller

import (
	"book-api/model"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var bookCollection *mongo.Collection = model.OpenCollection(model.Client, "book")

func BookGet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	var book []bson.M
	
	cursor, err := bookCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = cursor.All(ctx, &book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, gin.H{"data": book})
}

func BookByIdGet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var book bson.M
	
	err := bookCollection.FindOne(ctx, bson.M{"_id":id}).Decode(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer cancel()
	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
}

func BookPost(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	var book model.Book

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	validationErr := validate.Struct(book)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}
	result, insertErr := bookCollection.InsertOne(ctx, book)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}

func BookDelete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	err := bookCollection.FindOneAndDelete(ctx, bson.M{"_id": id}).Err()
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, gin.H{
		"message": "book delete",
	})
}

func BookUpate(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)

	var book model.Book
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedBook, err := bookCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set",
			bson.M{
				"title": book.Title,
				"author": book.Author,
				"year": book.Year,
				},
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, gin.H{
		"data": updatedBook,
	})
}