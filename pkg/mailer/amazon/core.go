/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package amazon

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/adrianpk/poslan/internal/config"
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Mailer is an amazon Mailer.
type Mailer struct {
	ctx    context.Context
	cfg    *config.Config
	logger log.Logger
	client *ses.SES
}

// Init amazon SES mail server handler.
func Init(ctx context.Context, cfg *config.Config, log log.Logger) (*Mailer, error) {
	// Initialize a new mailer.
	return newMailer(ctx, cfg, log)
}

// Send an email.
func (m *Mailer) Send(em model.Email) (resend bool, err error) {

	email := newSESEmail(em.From, em.To, em.CC, em.BCC, em.Subject, em.Body, em.Charset)

	result, err := m.client.SendEmail(email)

	// Actually, all error cases are solved in the same way.
	// In case that, eventually, it is not required to modify
	// this behavior for some particular case, the following block
	// could be replaced by a single line of code:
	// return true, err
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {

			case ses.ErrCodeMessageRejected:
				// SES mail sending not succeed
				// It probably does not exist but we can try again
				return true, fmt.Errorf("cannot send the email: %s", err.Error())

			case ses.ErrCodeMailFromDomainNotVerifiedException:
				// SES cannot read MX record.
				// It probably does not exist but we can try again
				// just in cae it was a temporary failure
				return true, fmt.Errorf("target domain not verified: %s", err.Error())

			case ses.ErrCodeConfigurationSetDoesNotExistException:
				// Configuration error, try a resend.
				return true, fmt.Errorf("configuration error: %s", err.Error())

			default:
				// Default condition for SES related errors.
				return true, fmt.Errorf("cannot send the email: %s", err.Error())
			}
		} else {
			// Default condition for SES non codified errors.
			return true, fmt.Errorf("cannot send the email: %s", err.Error())
		}
		// Default condition for non SES codified errors.
		return true, fmt.Errorf("cannot send the email: %s", err.Error())
	}

	m.logger.Log(
		"level", config.LogLevel.Info,
		"package", "amazon",
		"method", "Send",
		"result", result.GoString(),
	)

	return false, nil
}

func newSESEmail(from, to, cc, bcc, subject, body, charset string) *ses.SendEmailInput {
	// Assemble the email.
	email := &ses.SendEmailInput{
		Destination: &ses.Destination{
			BccAddresses: []*string{
				aws.String(bcc),
			},
			CcAddresses: []*string{
				aws.String(cc),
			},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	}

	return email
}

func newMailer(ctx context.Context, cfg *config.Config, logger log.Logger) (*Mailer, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		return nil, err
	}

	// Create a SES session.
	clt := ses.New(sess)

	return &Mailer{
		ctx:    ctx,
		cfg:    cfg,
		logger: logger,
		client: clt,
	}, nil
}

// Start the mailer.
func (m *Mailer) Start() error {
	return nil
}

// Stop the mailers.
func (m *Mailer) Stop() error {
	return nil
}

// IsReady return true if mailer is ready.
func (m *Mailer) IsReady() bool {
	return true
}
