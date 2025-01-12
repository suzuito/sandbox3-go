package main

import (
	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/web"
)

func main() {
	w := web.New()

	e := gin.New()
	e.Use(gin.Recovery())
	w.SetEngine(e)

	e.Run()
}
