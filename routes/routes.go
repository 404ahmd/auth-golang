package routes

import (
    "auth-golang-jwt/handlers"
    "auth-golang-jwt/middleware"

    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    api := r.Group("/api")

    // Public routes
    auth := api.Group("/auth")
    {
        auth.POST("/register", handlers.Register)
        auth.POST("/login", handlers.Login)
    }

    // Protected routes
    protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", handlers.Profile)
    }
}