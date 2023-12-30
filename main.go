package main

import (
	"fmt"
	"os"

	rt "github.com/Kotzyk/go-short/api/route"
	"github.com/Kotzyk/go-short/helpers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	// Initialize Gin engine
	server := gin.Default()

	helpers.SetupPrometheus("/metrics", server)

	rt.SetUrlsRouter(server)
	// Start HTTP server
	port := os.Getenv("PORT")
	fmt.Printf("Port read as: %s\n", port)
	server.Run(fmt.Sprintf(":%s", port))
}
