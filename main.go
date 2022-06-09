package main

import (
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/zachary-cauchi/golang-angular-sample-app/internal/handlers"
)

func main() {
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := filepath.Ext(file)

		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}
	})

	r.GET("/message", handlers.GetMessageListHandler)
	r.POST("/message", handlers.AddMessageHandler)
	r.DELETE("/message/:id", handlers.DeleteMessageHandler)

	err := r.Run(":3000")

	if err != nil {
		panic(err)
	}
}
