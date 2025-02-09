package main

import (
	"github.com/legenhand/recipe-test/src/config"
	"github.com/legenhand/recipe-test/src/db"
	"github.com/legenhand/recipe-test/src/router"
	"github.com/legenhand/recipe-test/src/seeder"
	"log"
	"os"
)

func main() {
	config.LoadConfig()
	db.DB = db.SetupGormDB()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		seeder.RunSeeds(db.DB)
		return
	}

	routes := router.SetupRouter()

	if err := db.Migrate(db.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("Server is running on port %s", config.Cfg.ServerPort)
	if err := routes.Run(":" + config.Cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
