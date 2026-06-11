package handler

import (
	"net/http"

	"shortURL/service"

	"github.com/gin-gonic/gin"
)

type shortenReq struct {
	LongURL string `json:"long_url" binding:"required"`
}

func ShortenHandler(domain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req shortenReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := service.Shorten(req.LongURL, domain)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
