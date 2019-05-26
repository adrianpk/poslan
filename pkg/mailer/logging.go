/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianpk/poslan/internal/config"
	c "github.com/adrianpk/poslan/internal/config"
	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

type loggingMiddleware struct {
	ctx    context.Context
	cfg    *config.Config
	logger log.Logger
	next   Service
}

// SignIn is a logging middleware wrapper over another interface implementation of SignIn.
func (mw loggingMiddleware) SignIn(clientID, secret string) (output string, err error) {
	defer func(begin time.Time) {
		input := fmt.Sprintf("{%s, %s}", clientID, secret)
		mw.logger.Log(
			"level", c.LogLevel.Info,
			"method", "SignIn",
			"input", input,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.SignIn(clientID, secret)
	return
}

// SignOut is a logging middleware wrapper over another interface implementation of SignOut.
func (mw loggingMiddleware) SignOut(id uuid.UUID, token string) (err error) {
	defer func(begin time.Time) {
		input := fmt.Sprintf("{%s}", id.String())
		mw.logger.Log(
			"level", c.LogLevel.Info,
			"method", "SignOut",
			"input", input,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.SignOut(id, token)
	return
}

// Send is a logging middleware wrapper over another interface implementation of Send.
func (mw loggingMiddleware) Send(to, cc, bcc, subject, body, token string) (err error) {
	defer func(begin time.Time) {
		input := fmt.Sprintf("{%s, %s, %s, %s}", to, cc, bcc, body)
		mw.logger.Log(
			"level", c.LogLevel.Info,
			"method", "Send",
			"input", input,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.Send(to, cc, bcc, subject, body, token)
	return
}

func (mw loggingMiddleware) Context() context.Context {
	return mw.ctx
}

// Config returns service config.
func (mw loggingMiddleware) Config() *config.Config {
	return mw.cfg
}

// Remove is a logging middleware wrapper over another interface implementation of Remove.
func (mw loggingMiddleware) Logger() log.Logger {
	return mw.logger
}
