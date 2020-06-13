package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service struct {
	secret []byte
}

func New(secret string) *Service {
	return &Service{secret: []byte(secret)}
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID   string `json:"user_id"`
	ClientIP string `json:"client_ip"`
}

func (s *Service) Generate(userID, clientIP string, expires time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expires).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "File Browser",
		},
		UserID:   userID,
		ClientIP: clientIP,
	})
	return t.SignedString(s.secret)
}

func (s *Service) Parse(ts, clientIP string) (userID string, err error) {
	var claims tokenClaims
	_, err = jwt.ParseWithClaims(ts, &claims, func(_ *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}

	if claims.ClientIP != clientIP {
		return "", errors.New("client ip mismatch")
	}
	return claims.UserID, nil
}
