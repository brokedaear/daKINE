// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package validator provides interfaces and functions for validating types.
package validator

import (
	"go.brokedaear.com/pkg/errors"
)

// Verifiable is a type that can be validated against constraints defined in
// the Validate method. A Value method also exists on Verifiable to return
// the underlying content that had been verified.
type Verifiable interface {
	Validate() error
	Value() any
}

// Check iterates through a list of Verifiable types, calling their Validate
// methods to check if the type respects their constraints. It joins multiple
// errors together when returning an error.
func Check(types ...Verifiable) error {
	if len(types) == 0 {
		return ErrNoTypesProvided
	}
	errs := make([]error, 0)
	for i, v := range types {
		if v == nil {
			errs = append(errs, ErrVerifiableIsNil)
		}
		err := v.Validate()
		if err != nil {
			errs = append(errs, errors.Errorf("%d. %v", i, err))
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

type VerifiableError string

func (e VerifiableError) Error() string {
	return string(e)
}

const (
	ErrNoTypesProvided VerifiableError = "no list of types provided"
	ErrVerifiableIsNil VerifiableError = "provided verifiable is nil"
)
