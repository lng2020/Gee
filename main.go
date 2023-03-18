package main

import (
	"gef"
	"net/http"
)

func main() {
	r := gef.New()
	r.GET("/", func(c *gef.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gef.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gef.Context) {
			// expect /hello?name=lng2020
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gef.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("/login", func(c *gef.Context) {
			c.JSON(http.StatusOK, gef.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":9999")
}
