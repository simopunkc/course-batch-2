package user

import (
	"course/internal/domain"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	privateKey []byte = []byte("mySignaturePrivateKey")
)

type UserUsecase struct {
	db *gorm.DB
}

func NewUserUsecase(db *gorm.DB) *UserUsecase {
	return &UserUsecase{
		db: db,
	}
}

func (uu UserUsecase) Register(c *gin.Context) {
	var user domain.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if user.Name == "" {
		c.JSON(400, gin.H{
			"message": "field name must required",
		})
		return
	}

	if user.Email == "" {
		c.JSON(400, gin.H{
			"message": "field email must required",
		})
		return
	}

	if user.Password == "" {
		c.JSON(400, gin.H{
			"message": "field password must required",
		})
		return
	}

	if len(user.Password) < 6 {
		c.JSON(400, gin.H{
			"message": "password must more than 5 character",
		})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	err = uu.db.Create(&user).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed when create user",
		})
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed generate token",
		})
		return
	}

	c.JSON(201, gin.H{
		"token": token,
	})
}

func (uu UserUsecase) Login(c *gin.Context) {
	var userRequest domain.User
	err := c.ShouldBind(&userRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}
	if userRequest.Email == "" || userRequest.Password == "" {
		c.JSON(400, gin.H{
			"message": "email/password salah",
		})
		return
	}

	var user domain.User
	err = uu.db.Where("email = ?", userRequest.Email).Take(&user).Error
	if err != nil || user.ID == 0 {
		c.JSON(400, gin.H{
			"message": "email/password salah",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "email/password salah",
		})
		return
	}
	token, _ := generateJWT(user.ID)
	c.JSON(200, gin.H{
		"token": token,
	})
}

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iss":     "edspert",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uu UserUsecase) DecriptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth invalid")
		}
		return privateKey, nil
	})

	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}

	if !parsedToken.Valid {
		return data, errors.New("invalid token")
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}
