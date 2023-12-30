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

	server.Run(fmt.Sprintf("0.0.0.0:%s", port))
}
