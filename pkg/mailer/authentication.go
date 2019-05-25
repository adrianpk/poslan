/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"github.com/adrianpk/poslan/pkg/auth"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type authenticationMiddleware struct {
	logger log.Logger
	auth   auth.SecServer
	next   Service
}

// SignIn is a logging middleware wrapper over another interface implementation of SignIn.
func (mw authenticationMiddleware) SignIn(username, password string) (output *model.User, err error) {
	authsrv := mw.auth
	authsrv.Authenticate(username, password)
	output, err = mw.next.SignIn(username, password)
	return
}

// SignOut is a logging middleware wrapper over another interface implementation of SignOut.
func (mw authenticationMiddleware) SignOut(id uuid.UUID) (err error) {
	err = mw.next.SignOut(id)
	return
}

// Send is a logging middleware wrapper over another interface implementation of Send.
func (mw authenticationMiddleware) Send(to, cc, bcc, subject, body string) (err error) {
	err = mw.next.Send(to, cc, bcc, subject, body)
	return
}

// Remove is a logging middleware wrapper over another interface implementation of Remove.
func (mw authenticationMiddleware) Logger() log.Logger {
	return mw.logger
}
