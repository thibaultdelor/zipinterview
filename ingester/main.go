package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ClickEvent is the json we received from the API
type ClickEvent struct {
	UserID int `json:"user_id" binding:"required"`
}

// IngestEvent take an event and put it in elastic
func IngestEvent(c *gin.Context) {
	// Validate input
	var input ClickEvent
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println(fmt.Sprintf("Invalid Json, error is \"%s\"", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("Error decoding Json")})
		return
	}

	// Add to elastic
	err := addElasticEvent(input.UserID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save event to Elastic :O"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": "event added!"})
	}
}

func main() {
	r := gin.Default()
	r.PUT("/", IngestEvent)
	r.Run()
}
