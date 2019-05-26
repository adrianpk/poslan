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
	"github.com/adrianpk/poslan/pkg/auth"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type authenticationMiddleware struct {
	ctx    context.Context
	cfg    *config.Config
	logger log.Logger
	auth   auth.SecServer
	next   Service
}

// SignIn is a logging middleware wrapper over another interface implementation of SignIn.
func (mw authenticationMiddleware) SignIn(clientID, secret string) (output string, err error) {
	output, err = mw.auth.Authenticate(clientID, secret)
	if err != nil {
		return "", err
	}
	return output, nil
}

// SignOut is a logging middleware wrapper over another interface implementation of SignOut.
func (mw authenticationMiddleware) SignOut(id uuid.UUID, token string) (err error) {
	err = mw.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	err = mw.next.SignOut(id, token)
	return
}

// Send is a logging middleware wrapper over another interface implementation of Send.
func (mw authenticationMiddleware) Send(to, cc, bcc, subject, body, token string) (err error) {
	err = mw.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	err = mw.next.Send(to, cc, bcc, subject, token, body)
	return
}

// Config returns service context.
func (mw authenticationMiddleware) Context() context.Context {
	return mw.ctx
}

// Config returns service config.
func (mw authenticationMiddleware) Config() *config.Config {
	return mw.cfg
}

// Config returns service logger.
func (mw authenticationMiddleware) Logger() log.Logger {
	return mw.logger
}
