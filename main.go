package main

import (
	"log"
	"net/http"

	"github.com/anangpermana/gin-gorm-base/controllers"
	"github.com/anangpermana/gin-gorm-base/initializers"
	"github.com/anangpermana/gin-gorm-base/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	if config.GinMode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Cors, config.ClientOrigin}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	server.GET("/", func(ctx *gin.Context) {
		message := "hello world"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "welcome go golang with gorm and postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	//routelist
	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
