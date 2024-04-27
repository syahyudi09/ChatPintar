package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syahyudi09/ChatPintar/App/utils/token"
)

// JWTMiddleware memvalidasi token JWT dan menambahkan klaim ke konteks Gin

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort() // Menghentikan eksekusi lebih lanjut
			return
		}

		// Menghilangkan prefiks "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Memvalidasi token
		claims, err := token.ValidateToken(tokenString) 
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Menyimpan klaim dalam konteks Gin
		c.Set("claims", claims)

		// Lanjutkan ke rute berikutnya
		c.Next()
	}
}