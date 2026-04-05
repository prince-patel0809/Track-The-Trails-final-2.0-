package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"track-the-trails/config"
	"track-the-trails/routes"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env only for local development
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
	config.ConnectDB()
	config.ConnectCloudinary()

	// auth routes
	routes.AuthRoutes()

	// todos routes
	routes.TodoRoutes()

	// projects routes
	routes.ProjectsRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on: http://localhost:" + port)

	http.ListenAndServe(":"+port, nil)
}
