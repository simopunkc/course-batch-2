package main

import (
	"course/internal/database"
	"course/internal/exercise"
	"course/internal/middleware"
	"course/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
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

	r.GET("/exercises/:id", middleware.WithJWT(userUsecase), exerciseService.GetExerciseByID)
	r.GET("/exercises/:id/score", middleware.WithJWT(userUsecase), exerciseService.CalculateUserScore)

	r.Run(":1234")
}
