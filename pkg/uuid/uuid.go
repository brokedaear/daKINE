// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package uuid wraps a UUID 7 generator.
package uuid

import "github.com/gofrs/uuid/v5"

// New generates a new UUID v7 and returns its string representation.
func New() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
