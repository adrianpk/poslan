/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import "github.com/adrianpk/poslan/pkg/model"

// Mailer interface
type Mailer interface {
	Send(model.Email) (resend bool, err error)
	Start() error
	Stop() error
	IsReady() bool
}
