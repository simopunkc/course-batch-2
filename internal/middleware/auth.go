package middleware

import (
	"context"
	"course/internal/user"
	"strings"

	"github.com/gin-gonic/gin"
)

func WithJWT(us *user.UserUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"message": "unauthorize",
			})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{
				"message": "unauthorize",
			})
			c.Abort()
		}

		auths := strings.Split(authHeader, " ")
		data, err := us.DecriptJWT(auths[1])
		if err != nil {
			c.JSON(401, gin.H{
				"message": "unauthorize",
			})
			c.Abort()
			return
		}
		ctxUserID := context.WithValue(c.Request.Context(), "user_id", data["user_id"])
		c.Request = c.Request.WithContext(ctxUserID)
		c.Next()
	}
}
