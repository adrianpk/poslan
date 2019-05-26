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

type contextKey string

// Request & response
// Sign in
type signInRequest struct {
	ClientID string `json:"clientID"`
	Secret   string `json:"secret"`
}

type signInResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"error,omitempty"`
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

func (c contextKey) String() string {
	return "poslan-" + string(c)
}
