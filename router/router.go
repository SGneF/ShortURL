package router

import (
	"shortURL/handler"

	"github.com/gin-gonic/gin"
)

func Setup(domain string) *gin.Engine {
	r := gin.Default()

	r.POST("/api/shorten", handler.ShortenHandler(domain))
	r.GET("/:shortCode", handler.RedirectHandler())

	return r
}
