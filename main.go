package main

import (
	"gef"
	"net/http"
)

func main() {
	r := gef.Default()
	r.GET("/", func(c *gef.Context) {
		c.String(http.StatusOK, "Hello lng2020\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gef.Context) {
		names := []string{"lng2020"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
}
