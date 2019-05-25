/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

// Service provides authentication and authorization services
type Service interface {
	SignIn(username, password string) (*model.User, error)
	SignOut(id uuid.UUID) error
	Send(to, cc, bcc, subject, body string) error
	Logger() log.Logger
}

// Mailer interface
type Mailer interface {
	Send(model.Email) (resend bool, err error)
	Start() error
	Stop() error
	IsReady() bool
}
