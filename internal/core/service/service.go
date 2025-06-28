// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package service defines services that can be used by controllers.
package service

type Service struct{}

func NewServices() *Service {
	return &Service{}
}
