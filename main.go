package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const PathWeb = "./web/dist/"

func main() {
	/*
	 * Import Environment
	 */
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set mode: release or debug
	gin.SetMode(os.Getenv("GIN_MODE"))

	/*
	 * Route
	 */
	route := gin.Default()

	// Share static files
	route.Static("/web", PathWeb)
	route.LoadHTMLGlob(PathWeb + "index.html")

	route.NoRoute(func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// use pusher.Push() to do server push
			if err := pusher.Push("/web/main.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}

		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Start Server
	err = route.RunTLS(":8080", os.Getenv("CERT_FILE_PATH"), os.Getenv("CERT_KEY_PATH"))
	if err != nil {
		log.Fatal("Error RunTLS server")
	}
}
