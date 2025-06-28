// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package server declares server related dependencies and functionalities.
package server

import (
	"net"
)

// Base is the base setup for any kind of server.
type Base struct {
	logger   Logger
	config   *config
	listener net.Listener
}

// NewBase validates and initializes the config, logger, and listener.
func NewBase(logger Logger, configOpts ...ConfigOpts) (*Base, error) {
	if logger == nil {
		return nil, ErrNilLogger
	}
	cfg, err := newConfig(configOpts...)
	if err != nil {
		return nil, err
	}
	var lis net.Listener
	var address string
	// If the socket path is set, use that instead of TCP.
	if len(cfg.socketPath.String()) > 0 {
		address = cfg.socketPath.String()
		lis, err = net.Listen("unix", address)
		if err != nil {
			return nil, err
		}
	} else {
		address, err = cfg.newURIAddress()
		if err != nil {
			return nil, err
		}
		lis, err = net.Listen("tcp", address)
		if err != nil {
			return nil, err
		}
	}
	return &Base{
		logger:   logger,
		config:   cfg,
		listener: lis,
	}, nil
}

type BaseError string

func (b BaseError) Error() string {
	return string(b)
}

var (
	ErrNilConfig BaseError = "config is nil"
	ErrNilLogger BaseError = "logger is nil"
)
