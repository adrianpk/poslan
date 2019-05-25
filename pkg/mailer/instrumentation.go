package mailer

import (
	"fmt"
	"time"

	// "github.com/go-kit/kit/log"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/google/uuid"
)

type instrumentationMiddleware struct {
	logger         log.Logger
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	next           Service
}

// SignIn is an instrumentation middleware wrapper over another interface implementation of SignIn.
func (mw instrumentationMiddleware) SignIn(username, password string) (output *model.User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SignIn", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.SignIn(username, password)
}

// SignOut is an instrumentation middleware wrapper over another interface implementation of SignOut.
func (mw instrumentationMiddleware) SignOut(id uuid.UUID) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "SignOut", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.SignOut(id)
}

// Send is an instrumentation middleware wrapper over another interface implementation of Send.
func (mw instrumentationMiddleware) Send(to, cc, bcc, subject, body string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Send", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.Send(to, cc, bcc, subject, body)
}

func (mw instrumentationMiddleware) Logger() log.Logger {
	return mw.logger
}
