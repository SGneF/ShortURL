package handler

import (
	"net/http"

	"shortURL/service"

	"github.com/gin-gonic/gin"
)

func RedirectHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		longURL, err := service.GetLongURL(c.Request.Context(), shortCode)
		if err != nil {
			c.String(http.StatusNotFound, "short url not found")
			return
		}

		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Redirect(http.StatusFound, longURL)
	}
}
