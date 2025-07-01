// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package domain

import "time"

type UserSession struct {
	Token     string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewUserSession(userID string, validFor time.Duration) *UserSession {
	now := time.Now().UTC()
	d := now.Add(validFor)
	return &UserSession{
		Token:     "",
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: d,
	}
}
