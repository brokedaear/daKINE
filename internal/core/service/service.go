// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package service defines services that can be used by controllers.
package service

import (
	"go.brokedaear.com/internal/common/telemetry"
	"go.brokedaear.com/internal/common/utils/loggers"
)

type Service struct{}

func NewServices() *Service {
	return &Service{}
}

// ServiceBase is a base type for all services.
type ServiceBase struct {
	logger loggers.Logger
	tel    telemetry.Telemetry
}

func NewServiceBase(logger loggers.Logger, tel telemetry.Telemetry) *ServiceBase {
	return &ServiceBase{
		logger: logger,
		tel:    tel,
	}
}
