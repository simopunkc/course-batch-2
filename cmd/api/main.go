package main

import (
	"course/internal/database"
	"course/internal/exercise"
	"course/internal/middleware"
	"course/internal/user"
	"log"
	"net/http"

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
	userUsecase := user.NewUserUsecase(db)
	r.POST("/register", userUsecase.Register)
	r.POST("/login", userUsecase.Login)
	r.POST("/exercises/", middleware.WithJWT(userUsecase), exerciseService.CreateExercise)
	r.POST("/questions/", middleware.WithJWT(userUsecase), exerciseService.CreateQuestion)

	r.GET("/exercises/:id", middleware.WithJWT(userUsecase), exerciseService.GetExerciseByID)
	r.GET("/exercises/:id/score", middleware.WithJWT(userUsecase), exerciseService.CalculateUserScore)

	r.Run(":1234")
}
