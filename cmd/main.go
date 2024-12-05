package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nietxan/LDM-Backend/web"
)

func main() {
	r := gin.Default()

	var err error

	err = web.InitDatabase()
	if err != nil {
		log.Fatal("Error loading database:", err)
	}
	
	err = web.InitEnv()
	if err != nil {
		log.Fatal("Error loading environmental variables:", err)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	r.POST("/singup", web.SingUp)
	r.POST("/singin", web.SingIn)
	r.GET("/update", web.OrderUpdateSocket)

	if err := r.Run(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
