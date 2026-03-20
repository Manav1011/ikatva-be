package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/manav1011/ikatva-be/internal/config"
	"github.com/manav1011/ikatva-be/internal/db"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	fmt.Printf("Config: %+v\n", cfg)
	// Load db
	pool := db.NewDB(cfg.DBSource)
	fmt.Println("DB Pool:", pool)
	// Init server
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("🚀 Server running on port", cfg.ServerPort)

	err := r.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	// Close db connection when the server shuts down
	defer pool.Close()
}
