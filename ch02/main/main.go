package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Person .
type Person struct {
	ID   int    `uri:"id" binding:"required" json:"id"`
	Name string `uri:"name" binding:"required" json:"name"`
	Book string `form:"book" json:"book"`
}

func main() {
	router := gin.Default()

	router.GET("/person/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.Status(404)
			return
		}
		if err := c.ShouldBind(&person); err != nil {
			c.Status(404)
			return
		}
		// 这种连续绑定的它是按顺序来进行判断的，首先进行
		// 这个c.ShouldBindUri(&person)绑定的判断，它会拿person里面的每个字段去c里面找相应的值
		// （且这个值是要在它绑定的类型里面找，例如：Uri或form）。如果能在相应的反射里面找到的话，
		// 它就会进行赋值操作，否则它会继续看这个字段是否存在banding的属性是否存在，如果存在就进行
		// banding的操作，如果banding 的操作失败则报错，如果没有banding的话，它就不进行赋值操作，
		// 让下一个ShouldBind进行赋值，以此类推
		fmt.Printf("%#v", person)
		c.JSON(200, person)
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
