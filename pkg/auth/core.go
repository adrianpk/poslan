package auth

import (
	"context"
	"errors"
	"time"

	"github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

const (
	expSeconds = 240
)

// SecServer is an authentication service.
type SecServer interface {
	// Authenticate generate the bearer token.
	Authenticate(string, string) (string, error)
	// ValidateToken ensure that the authentication token is valid.
	ValidateToken(string) error
}

// Server is an omplementation of SecServer.
type Server struct {
	ctx      context.Context
	cfg      *config.Config
	Logger   log.Logger
	key      []byte
	clientDB map[string]string
	usersDB  map[string]*model.User
}

type customClaims struct {
	ClientID string `json:"clientID"`
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func generateToken(signingKey []byte, clientID string, user *model.User) (string, error) {
	claims := customClaims{
		clientID,
		user.ID.String(),
		user.Username,
		user.Name,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * expSeconds).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

// Authenticate a user.
func (s Server) Authenticate(clientID string, clientSecret string) (string, error) {
	if s.validSecret(clientID, clientSecret) {

		user := s.userBySecret(clientSecret)

		signed, err := generateToken(s.signingKey(), clientID, user)
		if err != nil {
			return "", errors.New("token generation error")
		}
		return signed, nil
	}
	return "", errors.New("wrong credentials")
}

// ValidateToken validate if token is valid.
func (s Server) ValidateToken(token string) error {

	t, err := jwt.Parse(token, s.Keys())

	if t.Valid {
		return nil
	}

	if v, ok := err.(*jwt.ValidationError); ok {

		if v.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token malformed
			return err

		} else if v.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token expired
			return err

		} else {
			// Other possible error
			return err
		}
	}
	return errors.New("invalid token")
}

// Keys returns a function used to generate the signing keys.
func (s Server) Keys() func(token *jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		return s.signingKey(), nil
	}
}

// This is a PoC, in a real world implementation
// this client database would be supported by
// some persistence mechanism.
// Clients map a client with its assigned secret.
func (s Server) clients() map[string]string {
	s.clientDB = make(map[string]string)
	s.clientDB["dd74cb9cfb5a4f1cac4d"] = "a5ee54c8a21a4c61820f88f14c30fa5b"
	s.clientDB["984fd4bdcb374aa7836a"] = "98d28599e5554a9ea4ada53feae924ff"
	return s.clientDB
}

func (s Server) validSecret(clientID, secret string) (valid bool) {
	clientsDB := s.clients()
	return clientsDB[clientID] == secret
}

// This is a PoC, in a real world implementation
// this user database would be supported by
// some persistence mechanism.
func (s Server) users() map[string]*model.User {
	s.usersDB = make(map[string]*model.User)
	iddp, _ := uuid.Parse("a4d50605-54b0-4c80-a7ef-60b6516998cc")
	s.usersDB["a5ee54c8a21a4c61820f88f14c30fa5b"] = &model.User{
		ID:       iddp,
		Username: "Diana Prince",
		Password: "3ae61c5e5af2276ee452237e573a8cf",
		Name:     "Diana",
		Email:    "diana.p@gmail.com",
	}
	idck, _ := uuid.Parse("c6da610a-d858-401d-8f0a-381cc6d6921a")
	s.usersDB["98d28599e5554a9ea4ada53feae924ff"] = &model.User{
		ID:       idck,
		Username: "Clark Kent",
		Password: "64cd0e8b4a00b7d22f40b124413ad16d",
		Name:     "Clarck",
		Email:    "clark.k@gmail.com",
	}
	return s.usersDB
}

func (s Server) signingKey() []byte {
	// signingkey: This a PoC.
	// For a production ready app use something like this to generate the keys
	// ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key
	// openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
	// cat jwtRS256.key
	// cat jwtRS256.key.pub
	// And of course get it from another place (i.e.: Hashicorp Vault)
	return []byte("55a124b9")
}

func (s Server) userBySecret(secret string) *model.User {
	usersDB := s.users()
	return usersDB[secret]
}

func (s Server) user(secret string) *model.User {
	usersDB := s.users()
	return usersDB[secret]
}

// Context returns service context.
func (s Server) Context() context.Context {
	return s.ctx
}

// Config returns service config.
func (s *Server) Config() *config.Config {
	return s.cfg
}

// UserData extract and return username and email from token.
func UserData(tokenString string) (userData map[string]string, err error) {
	userData = make(map[string]string)
	cs, err := GetClaims(tokenString)

	if err != nil {
		return nil, err
	}

	for key, val := range cs {
		if key == "clientID" {
			userData["clientID"] = val.(string)
		}
		if key == "userID" {
			userData["userID"] = val.(string)
		}
		if key == "Name" {
			userData["name"] = val.(string)
		}
		if key == "email" {
			userData["email"] = val.(string)
		}
	}

	return userData, nil
}

// GetClaims extract and return claims from token.
func GetClaims(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("55a124b9"), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
