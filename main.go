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

	r.GET("/hello", func(c *gef.Context) {
		// expect /hello?name=lng2020
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gef.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gef.Context) {
		c.JSON(http.StatusOK, gef.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
