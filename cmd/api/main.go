package main

import (
	"course/internal/database"
	"course/internal/exercise"
	"course/internal/middleware"
	"course/internal/user"
	"course/internal/user/repository"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})

	db := database.NewConnDatabase()
	exerciseService := exercise.NewExerciseUsecase(db)
	// repo := repository.NewDatabaseRepo(db)
	repo := repository.NewMicroserviceRepo()
	userUsecase := user.NewUserUsecase(repo)
	r.POST("/exercises", middleware.WithJWT(userUsecase), exerciseService.CreateExercise)
	r.POST("/exercises/:id/questions", middleware.WithJWT(userUsecase), exerciseService.CreateQuestion)
	r.POST("/exercises/:id/questions/:qid/answers", middleware.WithJWT(userUsecase), exerciseService.CreateAnswer)

	r.GET("/exercises/:id", middleware.WithJWT(userUsecase), exerciseService.GetExerciseByID)
	r.GET("/exercises/:id/score", middleware.WithJWT(userUsecase), exerciseService.CalculateUserScore)

	r.Run(":" + os.Getenv("APP_PORT"))
}
