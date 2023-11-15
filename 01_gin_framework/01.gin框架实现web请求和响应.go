package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/**
* Description:
* @Author Hollis
* @Create 2023-10-27 21:08
 */
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./index.html")

	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	err := r.Run(":9090")
	if err != nil {
		log.Println(err)
		return
	}
}
