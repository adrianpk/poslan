package auth

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"

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
	ctx      context.Context
	cfg      *config.Config
	logger   log.Logger
	key      []byte
	clientDB map[string]string
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
	if s.validSecret(clientID, clientSecret) {
		signed, err := generateToken(s.key, clientID)
		if err != nil {
			return "", errors.New("token generation error")
		}
		return signed, nil
	}
	return "", errors.New("wrong credentials")
}

// This is a PoC, in a real world implementation
// this client database would be supported by
// some persistence mechanism.
// Clients map a client with its assigned secret.
func (s Server) clients() map[string]string {
	s.clientDB = make(map[string]string)
	s.clientDB["clt1"] = "52b3d83e"
	s.clientDB["clt2"] = "1580c230"
	return s.clientDB
}

func (s Server) validSecret(clientID, secret string) (valid bool) {
	clientsDB := s.clients()
	return clientsDB[clientID] == secret
}

// Context returns service context.
func (s Server) Context() context.Context {
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
