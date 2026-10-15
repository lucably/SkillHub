package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func validateService(service Service) error {

	if strings.TrimSpace(service.Name) == "" {
		return errors.New("Name must not be empty")
	}

	if service.Price <= 0 {
		return errors.New("Price must be greater than zero")
	}

	return nil
}

func welcomeHandler(c *gin.Context) {
	response := map[string]string{
		"status":  "ok",
		"message": "Welcome, SkillHub!",
	}

	c.JSON(http.StatusOK, response)
}

func servicesHandler(c *gin.Context) {
	var service Service

	err := c.ShouldBindJSON(&service)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateService(service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func main() {
	router := gin.Default()
	router.GET("/welcome", welcomeHandler)
	router.POST("/services", servicesHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
