package router

import (
	"github.com/gin-gonic/gin"
	"github.com/legenhand/recipe-test/src/controller"
	"github.com/legenhand/recipe-test/src/middleware"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/submit-email", controller.SubmitEmail)
		authGroup.GET("/magic-link", controller.MagicLink)
	}

	protectedGroup := router.Group("/")
	protectedGroup.Use(middleware.Auth())
	{
		protectedGroup.GET("/inventory", controller.GetInventory)
		protectedGroup.POST("/recipe", controller.CreateRecipe)
	}
	return router
}
