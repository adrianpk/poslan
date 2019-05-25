/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package amazon

import (
	"context"

	"github.com/go-kit/kit/log"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/adrianpk/poslan/internal/config"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Mailer is an amazon Mailer.
type Mailer struct {
	ctx    context.Context
	cfg    *config.Config
	logger log.Logger
	client *ses.SES
}
