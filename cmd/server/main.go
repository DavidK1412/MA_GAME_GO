package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/org/ranas-bdi-backend/internal/adapters/db/postgres"
	httpadapter "github.com/org/ranas-bdi-backend/internal/adapters/http"
	"github.com/org/ranas-bdi-backend/internal/app/usecase"
	platformdb "github.com/org/ranas-bdi-backend/internal/platform/db"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()
	if err := platformdb.InitFromEnv(ctx); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	defer platformdb.Close()

	pool := platformdb.MustGet()

	sessionRepo := postgres.NewSessionRepository(pool)
	matchRepo := postgres.NewMatchRepository(pool)
	moveRepo := postgres.NewMoveRepository(pool)
	difficultyRepo := postgres.NewDifficultyRepository(pool)

	sessionService := usecase.NewSessionService(sessionRepo)
	matchService := usecase.NewMatchService(matchRepo)
	moveService := usecase.NewMoveService(moveRepo)
	difficultyService := usecase.NewDifficultyService(difficultyRepo)

	handler := httpadapter.NewHandler(sessionService, matchService, moveService, difficultyService)
	router := handler.Router()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
