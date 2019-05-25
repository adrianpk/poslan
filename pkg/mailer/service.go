/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
// jwtKeys *config.JWTKeys
)

// Service provides authentication and authorization services
type Service interface {
	SignIn(username, password string) (*model.User, error)
	SignOut(id uuid.UUID) error
	Send(to, cc, bcc, subject, body string) error
	Logger() log.Logger
}

type service struct {
	mux     sync.Mutex
	name    string
	ctx     context.Context
	cfg     *config.Config
	logger  log.Logger
	mailers []Mailer
}

// Interface implementation

// SignIn lets a user sign in providing username and password.
func (s service) SignIn(username, password string) (*model.User, error) {
	return &model.User{}, errors.New("not implemented")
}

// SignOut lets a user sign out.
func (s service) SignOut(id uuid.UUID) error {
	// TODO: Close session implementation.
	return errors.New("not implemented")
}

// Send lets the user send a mail.
func (s service) Send(to, cc, bcc, subject, body string) error {

	from := "fix:username" // TODO: Get from user in session data.
	m := makeEmail(from, to, cc, bcc, subject, body)
	s.logger.Log("message", fmt.Sprintf("+v", m))

	return errors.New("not implemented:")
}

// Misc

// StarMailers is used in service startup
// to start each configured mailer.
func (s *service) StartMailers() {
	for _, m := range s.mailers {
		m.Stop()
	}
}

// StarMailers is used in service stop
// to stop each configured mailer.
func (s *service) StopMailers() {
	for _, m := range s.mailers {
		m.Start()
	}
}

// Logger returns service imterface implemention logger.
func (s service) Logger() log.Logger {
	return s.logger
}

// Utility functions
func makeEmail(from, to, cc, bcc, subject, body string) *model.Email {
	return &model.Email{
		ID:      uuid.New(),
		From:    from,
		To:      to,
		CC:      cc,
		BCC:     bcc,
		Subject: subject,
		Body:    body,
		Charset: charset,
	}
}

func passwordMatches(digest, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(digest), []byte(password))
	if err != nil {
		return false
	}
	return true
}