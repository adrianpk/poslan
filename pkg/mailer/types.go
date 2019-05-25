/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package mailer

import (
	"github.com/adrianpk/poslan/pkg/model"
	"github.com/google/uuid"
)

// Request & response

// Sign in
type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	User *model.User `json:"user,omitempty"`
	Err  string      `json:"error,omitempty"`
}

// Sign out
type signOutRequest struct {
	ID uuid.UUID `json:"id,omitempty"`
}

type signOutResponse struct {
	Err string `json:"error,omitempty"`
}

// Send
type sendRequest struct {
	To      string `json:"to,omitempty"`
	Cc      string `json:"cc,omitempty"`
	Bcc     string `json:"bcc,omitempty"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
}

type sendResponse struct {
	Email *model.Email `json:"email,omitempty"`
	Err   string       `json:"error,omitempty"`
}
