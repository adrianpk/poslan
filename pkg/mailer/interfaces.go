/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"context"

	"github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

// Service provides authentication and authorization services
type Service interface {
	Context() context.Context
	Config() *config.Config
	Logger() log.Logger
	SignIn(clientID, secret string) (string, error)
	SignOut(id uuid.UUID, token string) error
	Send(to, cc, bcc, subject, body, token string) error
}

// Mailer interface
type Mailer interface {
	Send(model.Email) (resend bool, err error)
	Start() error
	Stop() error
	IsReady() bool
}
