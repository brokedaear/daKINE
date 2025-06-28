// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package loggers defines various loggers that can log things.
package loggers

// Logger describes a custom logging implementation. For example, one might
// implement a library like logrus, Zap, Charm, or use Go's slog library
// as their preferred logger. This interface abstracts that functionality.
//
// The methods on Logger are kept intentionally minimal. This makes the caller
// think less when using Logger (that's good for everyone involved) and gives
// the caller the option to format their own string, since there is no <level>f
// method on Logger (such as "Infof" or "Errorf").
type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Sync() error
	Close() error
}
