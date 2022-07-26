package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// import 的是路径，用的是package

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	fmt.Println(123)
}
