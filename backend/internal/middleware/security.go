package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CreateAccessToken(userID int64, secret string, expireHours int) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(time.Duration(expireHours) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// DecodeAccessToken returns the user ID encoded in the token, or ok=false if
// the token is missing, malformed, or expired.
func DecodeAccessToken(tokenString, secret string) (userID int64, ok bool) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, false
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	sub, ok := claims["sub"]
	if !ok {
		return 0, false
	}
	switch v := sub.(type) {
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}
