package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/PrasadNaik1310/driveClone/db"
	"github.com/PrasadNaik1310/driveClone/handlers"
	"github.com/PrasadNaik1310/driveClone/middleware"

	//"github.com/PrasadNaik1310/driveClone/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
		log.Fatal("env file not found, main file error")
		return
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
		log.Print(err)
		log.Fatal("Error in initializing db . coming from main...")

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
		{
			auth.POST("/login", handlers.Login)
		}
		protected := api.Group("", middleware.AuthMiddleWare())
		folder := protected.Group("/folder")
		{
			folder.POST("/CreateFolder", handlers.CreateFolder)
			folder.GET("/DeleteFolder", handlers.DeleteFolder)
		}
		file := protected.Group("/file")
		{
			file.POST("/UploadFile", handlers.UploadFile)
			file.GET("/DownloadFile/:id", handlers.DownloadFile)
			file.DELETE("/DeleteFile/:id", handlers.DeleteFile)
		}
	}
	port := os.Getenv("port")
	if port == "" {
		log.Println("Port not found , going fo rport 8080")
		port = "8080"
	}
	go func() {
		srv := &http.Server{
			Addr:    ":" + port,
			Handler: r,
		}
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Signal recived for shutdown , going for shutdown")

}
