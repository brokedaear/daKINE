// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package domain_test

import (
	"testing"

	"go.brokedaear.com/app/domain"
	"go.brokedaear.com/pkg/assert"
	"go.brokedaear.com/pkg/test"
)

type EnvironmentValidateTestCase struct {
	test.CaseBase
	Env domain.Environment
}

func TestEnvironment_Validate(t *testing.T) {
	tests := []EnvironmentValidateTestCase{
		{
			CaseBase: test.NewCaseBase("valid development environment", nil, false),
			Env:      domain.EnvDevelopment,
		},
		{
			CaseBase: test.NewCaseBase("valid staging environment", nil, false),
			Env:      domain.EnvStaging,
		},
		{
			CaseBase: test.NewCaseBase("valid production environment", nil, false),
			Env:      domain.EnvProduction,
		},
		{
			CaseBase: test.NewCaseBase("empty environment", "environment is nil", true),
			Env:      domain.Environment{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				err := tt.Env.Validate()
				assert.ErrorOrNoError(t, err, tt.WantErr)
				if tt.WantErr && err != nil {
					assert.Equal(t, err.Error(), tt.Want.(string))
				}
			},
		)
	}
}

type EnvironmentStringTestCase struct {
	test.CaseBase
	Env domain.Environment
}

func TestEnvironment_String(t *testing.T) {
	tests := []EnvironmentStringTestCase{
		{
			CaseBase: test.NewCaseBase("development environment string", "development", false),
			Env:      domain.EnvDevelopment,
		},
		{
			CaseBase: test.NewCaseBase("staging environment string", "staging", false),
			Env:      domain.EnvStaging,
		},
		{
			CaseBase: test.NewCaseBase("production environment string", "production", false),
			Env:      domain.EnvProduction,
		},
		{
			CaseBase: test.NewCaseBase("empty environment string", "", false),
			Env:      domain.Environment{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				assert.Equal(t, tt.Env.String(), tt.Want.(string))
			},
		)
	}
}

type EnvironmentValueTestCase struct {
	test.CaseBase
	Env domain.Environment
}

func TestEnvironment_Value(t *testing.T) {
	tests := []EnvironmentValueTestCase{
		{
			CaseBase: test.NewCaseBase("development environment value", "development", false),
			Env:      domain.EnvDevelopment,
		},
		{
			CaseBase: test.NewCaseBase("staging environment value", "staging", false),
			Env:      domain.EnvStaging,
		},
		{
			CaseBase: test.NewCaseBase("production environment value", "production", false),
			Env:      domain.EnvProduction,
		},
		{
			CaseBase: test.NewCaseBase("empty environment value", "", false),
			Env:      domain.Environment{},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				assert.Equal(t, tt.Env.Value().(string), tt.Want.(string))
			},
		)
	}
}

type EnvFromStringTestCase struct {
	test.CaseBase
	Input string
}

func TestEnvFromString(t *testing.T) {
	tests := []EnvFromStringTestCase{
		{
			CaseBase: test.NewCaseBase("development string", domain.EnvDevelopment, false),
			Input:    "development",
		},
		{
			CaseBase: test.NewCaseBase("staging string", domain.EnvStaging, false),
			Input:    "staging",
		},
		{
			CaseBase: test.NewCaseBase("production string", domain.EnvProduction, false),
			Input:    "production",
		},
		{
			CaseBase: test.NewCaseBase("empty string", "invalid environment", true),
			Input:    "",
		},
		{
			CaseBase: test.NewCaseBase("invalid string", "invalid environment", true),
			Input:    "unknown",
		},
		{
			CaseBase: test.NewCaseBase("case sensitive test", "invalid environment", true),
			Input:    "DEVELOPMENT",
		},
		{
			CaseBase: test.NewCaseBase("whitespace string", "invalid environment", true),
			Input:    " development ",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.Name, func(t *testing.T) {
				got, err := domain.EnvFromString(tt.Input)
				assert.ErrorOrNoError(t, err, tt.WantErr)
				if tt.WantErr {
					assert.Equal(t, err.Error(), tt.Want.(string))
				} else {
					assert.Equal(t, got, tt.Want.(domain.Environment))
				}
			},
		)
	}
}

func TestEnvironment_GlobalVariables(t *testing.T) {
	envs := []struct {
		name string
		env  domain.Environment
		str  string
	}{
		{"EnvDevelopment", domain.EnvDevelopment, "development"},
		{"EnvStaging", domain.EnvStaging, "staging"},
		{"EnvProduction", domain.EnvProduction, "production"},
	}

	for _, e := range envs {
		t.Run(
			e.name, func(t *testing.T) {
				assert.Equal(t, e.env.String(), e.str)
				assert.Equal(t, e.env.Value().(string), e.str)
			},
		)
	}
}

func TestEnvironment_Consistency(t *testing.T) {
	envs := []domain.Environment{domain.EnvDevelopment, domain.EnvStaging, domain.EnvProduction}

	for _, env := range envs {
		t.Run(
			"consistency_"+env.String(), func(t *testing.T) {
				reconstructed, err := domain.EnvFromString(env.String())
				assert.NoError(t, err)
				assert.Equal(t, reconstructed, env)
			},
		)
	}
}

func TestEnvironment_ValidateInvalid(t *testing.T) {
	// Test invalid environment validation by getting an invalid env from EnvFromString
	invalidEnv, _ := domain.EnvFromString("unknown")
	err := invalidEnv.Validate()
	assert.ErrorOrNoError(t, err, true)
	assert.Equal(t, err.Error(), "invalid environment")
}
