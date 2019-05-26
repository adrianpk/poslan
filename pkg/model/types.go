/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package model

import (
	"github.com/google/uuid"
)

// User model
type User struct {
	ID             uuid.UUID
	Username       string
	Password       string
	PasswordDigest string
}

// Email model
// A non PoC implementations should accespt lists, slices,
// for to, cc, and bcc fields
type Email struct {
	ID      uuid.UUID
	Name    string
	From    string
	To      string
	CC      string
	BCC     string
	Subject string
	Body    string
	Charset string
}
