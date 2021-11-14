package minigin

func main() {
	e := New()
	e.GET("/", func(c *Context) {
		c.Writer.Write([]byte("hello"))
	})
	e.Run(":8000")
}
