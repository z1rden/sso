package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"sso/internal/sso/auth/models"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
