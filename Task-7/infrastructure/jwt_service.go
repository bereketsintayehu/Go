package infrastructure

import (
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type JWTService struct {
	secret []byte
}

func NewJWTService() *JWTService {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		panic("JWT_SECRET environment variable not set")
	}
	

	return &JWTService{secret: secret}
}

func (j *JWTService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
}

func (j *JWTService) ExtractClaims(token *jwt.Token) (map[string]interface{}, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}

func (j *JWTService) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}