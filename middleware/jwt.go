package middleware
import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct{
	UserId uint `json:"user_id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uint, email string)(string, error){
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "U3Qskye4dI8dm0yJC0lvNrDNYGz2r0SEy2xMuuzbmhH"
	}

	claims := Claims{
		UserId : userId,
		Email : email,
		RegisteredClaims : jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header missing or invalid",
            })
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        secret := os.Getenv("JWT_SECRET")
        if secret == "" {
            secret = "your-secret-key"
        }

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            return
        }

        // Set user info ke context
        c.Set("user_id", claims.UserId)
        c.Set("email", claims.Email)
        c.Next()
    }
}
