package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/manav1011/ikatva-be/internal/config"
	"github.com/manav1011/ikatva-be/internal/db"
	sqldb "github.com/manav1011/ikatva-be/internal/db/sqlc"
	"github.com/manav1011/ikatva-be/internal/user"
	"github.com/manav1011/ikatva-be/internal/user/handler"
	"github.com/manav1011/ikatva-be/internal/user/repository"
	"github.com/manav1011/ikatva-be/internal/user/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/manav1011/ikatva-be/docs" // swagger docs (swag init)
)

// @title           Ikatva API
// @version         1.0
// @description     HTTP API for Ikatva backend.
// @host            localhost:8080
// @BasePath        /v1
// @schemes         http

// health returns service availability.
// @Summary      Health check
// @Description  Returns ok if the service is running.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func main() {
	cfg := config.LoadConfig()

	pool := db.NewDB(cfg.DBSource)

	sqlDB := stdlib.OpenDBFromPool(pool)
	defer sqlDB.Close()

	queries := sqldb.New(sqlDB)
	userRepo := repository.NewUserRepository(queries)
	userSvc := service.NewUserService(userRepo, cfg)
	userHandler := handler.NewUserHandler(userSvc)

	r := gin.Default()

	// Swagger UI: open /swagger/index.html (Gin cannot mix /swagger + /swagger/*any on the same prefix).
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	v1.GET("/health", health)

	user.RegisterRoutes(v1, userHandler)

	log.Println("server running on port", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
