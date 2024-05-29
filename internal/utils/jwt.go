package jwt

import (
	g "backend/internal/global"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var (
	ErrTokenExpired     = errors.New("token 已过期，请重新登录")
	ErrTokenNotValidYet = errors.New("token 无效，请重新登录")
	ErrTokenMalformed   = errors.New("token 不正确，请重新登录")
	ErrTokenNotValid    = errors.New("这不是一个 token ，请重新登录")
)

type MyClaim struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func GenToken(userID int, username, email string) (string, error) {
	secret := []byte(g.Conf.JWT.Secret)
	expireHour := g.Conf.JWT.Expire
	claim := &MyClaim{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    g.Conf.JWT.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString(secret)
}
func ParseToken(tokenString string) (*MyClaim, error) {
	secret := []byte(g.Conf.JWT.Secret)
	jwtToken, err := jwt.ParseWithClaims(
		tokenString, &MyClaim{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
	if err != nil {
		switch ve, ok := err.(*jwt.ValidationError); ok {
		case ve.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, ErrTokenMalformed
		case ve.Errors&jwt.ValidationErrorExpired != 0:
			return nil, ErrTokenExpired
		case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, ErrTokenNotValidYet
		}
	}
	if claims, ok := jwtToken.Claims.(*MyClaim); ok && jwtToken.Valid {
		return claims, nil
	}
	return nil, ErrTokenNotValid
}
