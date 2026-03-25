package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PrasadNaik1310/driveClone/db"
	"github.com/PrasadNaik1310/driveClone/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {

	var user struct {
		UserEmail string `json:"useremail"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var loginUser models.User
	if err := db.DB.Where("useremail=?", loginUser.UserEmail).First(&loginUser.UserEmail).Error; err != nil {
		log.Printf("User not found %s. Error from loginHandler...", loginUser.UserEmail)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return

	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Printf("jwt secret not found in environment , check for secret config ")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "ACCESS TOKEN NOT CONFIGURED , INTERNAL SERVER SIDE ERROR ."})
		return
	}

	log.Printf("JWT_SECRET for token generation: %s", jwtSecret[:10]+"...")

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userEmail": user.UserEmail,
		//"user":      user.userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("Token generated successfully for user: %s (ID: %d)", loginUser.UserEmail, loginUser.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"userEmail": loginUser.UserEmail,
			"userId":    loginUser.UserId,
			//"phone":      patient.Phoneno,
		},
	})
}
