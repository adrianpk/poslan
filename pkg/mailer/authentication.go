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
func (mw authenticationMiddleware) SignIn(ctx context.Context, clientID, secret string) (output string, err error) {
	output, err = mw.auth.Authenticate(clientID, secret)
	if err != nil {
		return "", err
	}
	return output, nil
}

// SignOut is a logging middleware wrapper over another interface implementation of SignOut.
func (mw authenticationMiddleware) SignOut(ctx context.Context, id uuid.UUID) (err error) {
	token, ok := AuthToken(ctx)
	if !ok {
		return errors.New("invalid token")
	}
	err = mw.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	err = mw.next.SignOut(ctx, id)
	return
}

// Send is a logging middleware wrapper over another interface implementation of Send.
func (mw authenticationMiddleware) Send(ctx context.Context, to, cc, bcc, subject, body string) (err error) {
	token, ok := AuthToken(ctx)
	if !ok {
		return errors.New("invalid token")
	}
	err = mw.auth.ValidateToken(token)
	if err != nil {
		return err
	}
	err = mw.next.Send(ctx, to, cc, bcc, subject, body)
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

// AuthToken gets the auth token from the context.
func AuthToken(ctx context.Context) (token string, ok bool) {
	token, ok = ctx.Value(authTokenCtxKey).(string)
	return token, ok
}
