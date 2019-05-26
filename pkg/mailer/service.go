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
	"github.com/adrianpk/poslan/internal/sys"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	mux       sync.Mutex
	name      string
	ctx       context.Context
	cfg       *config.Config
	logger    log.Logger
	providers []sys.Provider
}

// SignIn lets a user sign in providing username and password.
func (s *service) SignIn(clientID, secret string) (string, error) {
	return "", errors.New("not implemented")
}

// SignOut lets a user sign out.
func (s *service) SignOut(id uuid.UUID) error {
	// TODO: Close session implementation.
	return errors.New("not implemented")
}

// Send lets the user send a mail.
func (s *service) Send(to, cc, bcc, subject, body string) error {
	from := "fix:username" // TODO: Get from user in session data.
	m := makeEmail(from, to, cc, bcc, subject, body)
	s.logger.Log("message", fmt.Sprintf("%+v", m))

	return errors.New("not implemented")
}

// Providers returns service providers.
func (s *service) Providers() []sys.Provider {
	return s.providers
}

// Misc
// Context returns service context.
func (s *service) Context() context.Context {
	return s.ctx
}

// Config returns service config.
func (s *service) Config() *config.Config {
	return s.cfg
}

// Logger returns service imterface implemention logger.
func (s *service) Logger() log.Logger {
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
