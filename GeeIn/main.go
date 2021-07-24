package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()
	// r.GET("/", func(c *gee.Context) {
	// 	fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
	// })

	// r.GET("/hello", func(c *gee.Context) {
	// 	for k, v := range c.Req.Header {
	// 		fmt.Fprintf(c.Writer, "Header[%q] = %q\n", k, v)
	// 	}
	// })

	// r.GET("/hello", func(c *gee.Context) {
	// 	c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	// })

	// r.GET("/hello/:name", func(c *gee.Context) {
	// 	c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
	// })

	// r.GET("/assets/*filepath", func(c *gee.Context) {
	// 	c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	// })
	r.Use(gee.Logger())
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1> Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(gee.Logger())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
		})
		v2.GET("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	log.Println("Serve at http://localhost:9999")
	r.Run(":9999")
}
