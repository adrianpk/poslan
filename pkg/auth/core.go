package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	expSeconds = 240
)

// SecServer is an authentication service.
type SecServer interface {
	Authenticate(string, string) (string, error)
}

type Server struct {
	key     []byte
	clients map[string]string
}

type customClaims struct {
	ClientID string `json:"clientID"`
	jwt.StandardClaims
}

func generateToken(signingKey []byte, clientID string) (string, error) {
	claims := customClaims{
		clientID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * expSeconds).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// Authenticate a user
func (s Server) Authenticate(clientID string, clientSecret string) (string, error) {
	if s.clients[clientID] == clientSecret {
		signed, err := generateToken(s.key, clientID)
		if err != nil {
			return "", errors.New("token generation error")
		}
		return signed, nil
	}
	return "", errors.New("wrong credentials")
}
