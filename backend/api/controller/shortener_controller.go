package controller

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/Kotzyk/go-short/api/db"
	"github.com/Kotzyk/go-short/api/model"
	"github.com/Kotzyk/go-short/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

var validate *validator.Validate

func ShortenUrl(c *gin.Context) {
	var (
		request model.ShortenRequest
		err     error
		id      string
	)

	validate = validator.New(validator.WithRequiredStructEnabled())

	rc := db.CreateClient(0)
	defer rc.Close()
	//make sure the request body is parsable to desired JSON at all
	c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// make sure the JSON fits the other rules described in the model
	err = validate.Struct(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//It's more intuitive to give link expiration time in seconds, even if redis expiry is in microseconds
	request.Expiry *= time.Second

	if request.CustomShort == "" {
		id = helpers.EncodeBase62(rand.Uint64())
		val, _ := rc.Get(db.Ctx, id).Result()

		//if the ID is simply unlucky and already exists, roll again until it's a free one
		for val != "" {
			id = helpers.EncodeBase62(rand.Uint64())
			val, _ = rc.Get(db.Ctx, id).Result()
		}
	} else {
		id = request.CustomShort
		val, _ := rc.Get(db.Ctx, id).Result()
		if val != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "URL Custom short is already in use"})
			return
		}
	}

	err = rc.Set(db.Ctx, id, request.URL, request.Expiry).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to server"})
		return
	}

	response := model.ShortenResponse{
		URL:    request.URL,
		Short:  id,
		Expiry: request.Expiry,
	}

	c.JSON(http.StatusOK, response)
}

func ResolveUrl(c *gin.Context) {
	url := c.Param("slug")

	rc := db.CreateClient(0)
	defer rc.Close()

	destination, err := rc.Get(db.Ctx, url).Result()

	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found in db"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Redirect(301, destination)
}
