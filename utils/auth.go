package utils

import (
	"errors"
	"scheduler/internal/app/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenGenerator interface {
	GenerateToken(uint, string, string) (*models.JWT, error)
}

type JWTGenerator struct {
	Secret []byte
}

func NewJWTGenerator() TokenGenerator {
	var jwtKey = []byte(GetEnvKey("JWTKEY"))
	return &JWTGenerator{
		Secret: jwtKey,
	}
}

// GenerateToken : to Generate JWT Access Token
func (j *JWTGenerator) GenerateToken(id uint, email, role string) (*models.JWT, error) {
	token, err := j.getAccessToken(id, email, role)
	if err != nil {
		return nil, errors.New(ErrJWT)
	}
	var jwtToken models.JWT
	jwtToken.ID = id
	jwtToken.AccessToken = token
	return &jwtToken, nil
}

func (j *JWTGenerator) getAccessToken(id uint, email, role string) (string, error) {
	expirationTime := time.Now().Add(AccessTokenInterval * time.Minute)
	claims := &models.Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.Secret)
}
