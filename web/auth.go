package web

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}


func InitEnv() error {
	if os.Getenv("env") == "TEST" {
		os.Setenv("JWT_SECRET", "TESTING")
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			return err
		}
	}

	return nil
}

func SingUp(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	var existingUser User
	if err := DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
		return
	}

	user := User{Name: input.Name, Email: input.Email, Password: string(hashedPassword)}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating user"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token, "name": user.Name, "email": user.Email})
}

func SingIn(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	var user User
	if err := DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect password"})
		return
	}

	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "name": user.Name, "email": user.Email})
}
