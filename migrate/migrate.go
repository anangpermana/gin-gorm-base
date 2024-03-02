package main

import (
	"fmt"
	"log"

	"github.com/anangpermana/gin-gorm-base/initializers"
	"github.com/anangpermana/gin-gorm-base/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Member{})
	fmt.Println("? Migration complete")
}
