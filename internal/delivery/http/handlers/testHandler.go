package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"

)

func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message" : "test handler OK"})
}