package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Constants for token expiration and secret key
const (
	AccessTokenLifetime = 5 * time.Hour // Masa berlaku token akses
)

var secretKey = []byte(os.Getenv("API_SECRET"))
// GenerateToken menghasilkan token akses JWT dengan klaim tertentu
func GenerateToken(userID, phoneNumber string) (string, error) {
	expirationTime := time.Now().Add(AccessTokenLifetime).Unix()

	// Membuat klaim token dengan informasi yang relevan
	claims := jwt.MapClaims{
		"user_id":      userID,
		"phone_number": phoneNumber,
		"exp":          expirationTime, // Waktu kadaluarsa
	}

	// Membuat token JWT baru dengan metode penandatanganan dan klaim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan kunci rahasia
	return token.SignedString(secretKey)
}


func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode penandatanganan adalah HMAC-SHA256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		// Jika ada kesalahan atau token tidak valid, kembalikan kesalahan
		return nil, err
	}

	// Type assertion untuk klaim
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &claims, nil
	}

	// Jika klaim tidak bisa di-casting, kembalikan kesalahan
	return nil, jwt.ErrInvalidKey
}
