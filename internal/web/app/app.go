package app

import (
	"github.com/gin-gonic/gin"
)

func App() {
	r := gin.Default()
	r.GET("/view", viewHandler)
	r.GET("/recommend/", recommendHandler)
	r.LoadHTMLGlob("internal/web/template/*")

	err := r.Run()
	if err != nil {
		panic("Could not start web server: " + err.Error())
	}
}
