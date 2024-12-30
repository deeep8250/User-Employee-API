package main

import (
	"user-management-system/config"
	"user-management-system/controllers"
	"user-management-system/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	config.DbConnect()
	r := gin.Default()

	//publc routes
	r.POST("/login", controllers.Login)
	r.POST("/signin", controllers.Signin)

	//private routes
	r.GET("/filter", middleware.JWTMiddleware(), controllers.GetFilteredData)
	r.PUT("/update", middleware.JWTMiddleware(), controllers.Update)
	r.DELETE("/delete", middleware.JWTMiddleware(), controllers.Delete)

	r.Run(":8080")

}
