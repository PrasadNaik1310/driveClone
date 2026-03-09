package main

import (
	"log"

	"github.com/PrasadNaik1310/driveClone/db"
	"github.com/PrasadNaik1310/driveClone/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		log.Fatal("env file not found, main file error")
	}
	r := gin.Default()
	/*r.Use(middleware.ErrorHandler())
	r.Use(middleware.RequestLogging())*/
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") //change in prod
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,PATCH,UPDATE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			log.Println("Recieved preflight request , responding with 204...")
			c.AbortWithStatus(204)
			return
		}
		c.Next()

	})

	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
		log.Println("Error in initializing db . coming from main...")

	}
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Service": "driveClone",
			"Status":  "Healthy",
		})
	})
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		auth.POST("/login", handlers.Login)
		{

		}
	}
}
