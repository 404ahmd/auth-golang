package main

import (
    "log"
    "os"

    "auth-golang-jwt/config"
    "auth-golang-jwt/models"
    "auth-golang-jwt/routes"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env (hanya di development)
    if os.Getenv("APP_ENV") != "production" {
        if err := godotenv.Load(); err != nil {
            log.Println("No .env file found, using system environment")
        }
    }

    config.ConnectDB()
    config.DB.AutoMigrate(&models.User{}, &models.Employee{})

    r := gin.Default()
    routes.SetupRoutes(r)

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    r.Run(":" + port)
}