/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package sendgrid

import (
	"context"
	"errors"
	"fmt"

	"github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/go-kit/kit/log"
	sg "github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SGProvider is a mail delivery provider.
type SGProvider struct {
	ctx      context.Context
	cfg      *config.Config
	logger   log.Logger
	client   *sg.Client
	name     string
	priority int
}

// Init amazon SendGrid mail server handler.
func Init(ctx context.Context, cfg *config.Config, log log.Logger) (*SGProvider, error) {
	ok := cfg.HasProviderType(config.ProviderType.SendGrid)
	if !ok {
		return nil, errors.New("no config associated to an SesGrid provider")
	}
	return newProvider(ctx, cfg, log)
}

// Send an mail.
func (p *SGProvider) Send(em *model.Email) (resend bool, err error) {
	email := newSGEmail(em.From, em.To, em.CC, em.BCC, em.Subject, em.Body, em.Charset)

	res, err := p.client.Send(email)

	if err != nil {
		return true, err
	}
	// If no errores but response status code != accepted (202)
	if res.StatusCode != 202 {
		msg := fmt.Sprintf("cannot send email - status code: '%d'", res.StatusCode)
		return true, errors.New(msg)
	}

	return false, nil
}

func newSGEmail(from, to, cc, bcc, subject, body, charset string) *sgmail.SGMailV3 {
	// Assemble the mail.
	f := sgmail.NewEmail(from, from)
	s := subject
	t := sgmail.NewEmail(to, to)
	tb := body
	hb := fmt.Sprintf("<html><html><div>%s</div></body></html>", body)
	e := sgmail.NewSingleEmail(f, s, t, tb, hb)
	fmt.Printf("EMAIL is %+v\n", e)
	return e
}

func newProvider(ctx context.Context, cfg *config.Config, logger log.Logger) (*SGProvider, error) {
	// Create a SendGrid session.
	p, ok := cfg.Provider(config.ProviderType.SendGrid)
	if !ok {
		return nil, fmt.Errorf("no provider of type '%s' in config", config.ProviderType.SendGrid)
	}

	clt := sg.NewSendClient(p.APIKey)

	return &SGProvider{
		ctx:      ctx,
		cfg:      cfg,
		logger:   logger,
		client:   clt,
		name:     p.Name,
		priority: p.Priority,
	}, nil
}

// Name return the provider name.
func (p *SGProvider) Name() string {
	return p.name
}

// Priority return the provider priorit
func (p *SGProvider) Priority() int {
	return p.priority
}

// Start the mailer.
func (p *SGProvider) Start() error {
	return nil
}

// Stop the mailer.
func (p *SGProvider) Stop() error {
	return nil
}

// IsReady return true if mailer is ready.
func (p *SGProvider) IsReady() bool {
	return true
}

// Client return the provider client.
func (p *SGProvider) Client() interface{} {
	return p.client
}
