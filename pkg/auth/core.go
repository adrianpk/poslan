package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/adrianpk/poslan/internal/config"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	expSeconds = 240
)

// SecServer is an authentication service.
type SecServer interface {
	Authenticate(string, string) (string, error)
}

// Server is an omplementation of SecServer.
type Server struct {
	ctx     context.Context
	cfg     *config.Config
	logger  log.Logger
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

// Context returns service context.
func (s *Server) Context() context.Context {
	return s.ctx
}

// Config returns service config.
func (s *Server) Config() *config.Config {
	return s.cfg
}

// Logger returns service logger.
func (s *Server) Logger() log.Logger {
	return s.logger
}
