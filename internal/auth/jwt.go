package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int) (string, error) {
	expHours := 24
	if ev := os.Getenv("JWT_EXP_HOUR"); ev != "" {
		if v, err := strconv.Atoi(ev); err == nil && v > 0 {
			expHours = v
		}
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", fmt.Errorf("jwt не найден")
	}
	return token.SignedString(secret)
}
