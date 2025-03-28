package web

import (
	"net/http"
	"vibe-check/scanner"

	"github.com/gin-gonic/gin"
)

func StartServer(path string) {
	r := gin.Default()
	r.Static("/static", "./web/static")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})
	r.GET("/scan", func(c *gin.Context) {
		issues := scanner.ScanPath(path, false, "")
		c.JSON(http.StatusOK, issues)
	})
	r.Run(":4444")
}
