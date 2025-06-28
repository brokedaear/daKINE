// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"go.brokedaear.com/pkg/errors"
	"go.brokedaear.com/pkg/validator"
)

// Config is the configuration for telemetry.
type Config interface {
	ServiceName() string
	ServiceVersion() string
	ServiceID() string
	ExporterConfig() ExporterConfig
	validator.Verifiable
}

type config struct {
	serviceName    ServiceName
	serviceVersion ServiceVersion
	serviceID      ServiceID
	exporterConfig ExporterConfig
}

// ExporterConfig implements Config.
func (c *config) ExporterConfig() ExporterConfig {
	return c.exporterConfig
}

// ServiceID implements Config.
func (c *config) ServiceID() string {
	return c.serviceID.String()
}

// ServiceName implements Config.
func (c *config) ServiceName() string {
	return c.serviceName.String()
}

// ServiceVersion implements Config.
func (c *config) ServiceVersion() string {
	return c.serviceVersion.String()
}

func NewConfig(
	serviceName, serviceVersion, serviceID string,
	exporterConfig ExporterConfig,
) (Config, error) {
	cfg := newConfig(serviceName, serviceVersion, serviceID, exporterConfig)

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func newConfig(
	serviceName, serviceVersion, serviceID string,
	exporterConfig ExporterConfig,
) *config {
	return &config{
		serviceName:    ServiceName(serviceName),
		serviceVersion: ServiceVersion(serviceVersion),
		serviceID:      ServiceID(serviceID),
		exporterConfig: exporterConfig,
	}
}

func (c *config) Validate() error {
	var errs []error

	err := c.serviceName.Validate()
	if err != nil {
		errs = append(errs, errors.Wrap(err, ErrInvalidServiceName.Error()))
	}

	err = c.serviceVersion.Validate()
	if err != nil {
		errs = append(errs, errors.Wrap(err, ErrInvalidServiceVersion.Error()))
	}

	err = c.serviceID.Validate()
	if err != nil {
		errs = append(errs, errors.Wrap(err, ErrNoServiceID.Error()))
	}

	err = c.exporterConfig.Validate()
	if err != nil {
		errs = append(errs, errors.Wrap(err, ErrInvalidExporterConfig.Error()))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (c *config) Value() any {
	return c
}

type ServiceName string

func (s ServiceName) Validate() error {
	return validateServiceName(string(s))
}

func (s ServiceName) Value() any {
	return string(s)
}

func (s ServiceName) String() string {
	return string(s)
}

type ServiceVersion string

func (s ServiceVersion) Validate() error {
	if strings.TrimSpace(string(s)) == "" {
		return ErrNoServiceVersion
	}
	// Check maximum length
	if len(s) > serviceVersionLimit {
		return ErrServiceVersionTooLong
	}

	// Basic semantic version pattern check (flexible)
	// Allow versions like: 1.0.0, v1.0.0, 1.0.0-beta, 1.0.0+build.1, etc.
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-+_"
	for _, char := range s {
		if !strings.ContainsRune(validChars, char) {
			return errors.Wrapf(ErrServiceVersionInvalidChar, "char %s", string(char))
		}
	}
	return nil
}

func (s ServiceVersion) Value() any {
	return string(s)
}

func (s ServiceVersion) String() string {
	return string(s)
}

type ServiceID string

func (s ServiceID) Validate() error {
	if strings.TrimSpace(string(s)) == "" {
		return ErrNoServiceID
	}
	return nil
}

func (s ServiceID) Value() any {
	return string(s)
}

func (s ServiceID) String() string {
	return string(s)
}

// ExporterConfig holds configuration for an OTEL exporter.
type ExporterConfig struct {
	// Type defines the type of exporter. There are three options:
	// GRPC, HTTP, or file.
	Type ExporterType

	// Endpoint is the endpoint where OTEL data will be sent to. It takes the shape
	// of a hostname and port. HTTP(S) endpoints must include a protocol scheme.
	Endpoint ExporterEndpoint

	// Insecure defines whether the exporter will use a secure
	// means of communication, such as TLS.
	Insecure bool

	// Headers defines any HTTP headers that could be sent with the
	// request to the endpoint.
	Headers map[string]string
}

func NewExporterConfig(opts ...ExportConfigOption) ExporterConfig {
	rawURL := "localhost:4317"
	u, _ := url.Parse(rawURL)
	t := ExporterTypeGRPC
	e := &ExporterConfig{
		Type: t,
		Endpoint: ExporterEndpoint{
			URL: rawURL,
			url: u,
			t:   &t,
		},
		Insecure: true,
		Headers:  make(map[string]string),
	}
	for _, opt := range opts {
		opt(e)
	}
	return *e
}

type ExporterEndpoint struct {
	URL string
	url *url.URL
	t   *ExporterType
}

func newExporterEndpoint(url string) ExporterEndpoint {
	return ExporterEndpoint{
		URL: url,
		url: nil,
		t:   nil,
	}
}

func (e ExporterEndpoint) Validate() error {
	if e.url == nil {
		return errors.New("exporter URL must be set")
	}
	if e.t == nil {
		return errors.New("exporter type unspecified")
	}
	return nil
}

func (e ExporterEndpoint) Value() any {
	return e
}

type ExportConfigOption func(*ExporterConfig)

func WithType(typ ExporterType) ExportConfigOption {
	return func(c *ExporterConfig) {
		c.Type = typ
	}
}

func WithEndpoint(endpoint string) ExportConfigOption {
	return func(c *ExporterConfig) {
		c.Endpoint = newExporterEndpoint(endpoint)
	}
}

// TODO: implement security for GRPC
// func WithSecurity(security type like TLS) ExportConfigOption {}

func WithHeaders(headers map[string]string) ExportConfigOption {
	return func(c *ExporterConfig) {
		c.Headers = headers
	}
}

func (e ExporterConfig) Validate() error {
	var errs []error

	err := e.Type.Validate()
	if err != nil {
		errs = append(errs, err)
	}

	switch e.Type {
	case ExporterTypeGRPC:
		err = e.validateGRPCEndpoint()
		if err != nil {
			errs = append(errs, err)
		}
	case ExporterTypeHTTP:
		err = e.validateHTTPEndpoint()
		if err != nil {
			errs = append(errs, err)
		}
	case ExporterTypeStdout:
		err = e.validateStdoutEndpoint()
		if err != nil {
			errs = append(errs, err)
		}
	}

	err = e.validateHeaders()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (e ExporterConfig) validateGRPCEndpoint() error {
	if strings.TrimSpace(e.Endpoint.URL) == "" {
		return ErrEndpointRequired
	}

	if strings.Contains(e.Endpoint.URL, "://") {
		return ErrGRPCEndpointNoScheme
	}

	if !strings.Contains(e.Endpoint.URL, ":") {
		return ErrGRPCEndpointMissingPort
	}

	return nil
}

func (e ExporterConfig) validateHTTPEndpoint() error {
	if strings.TrimSpace(e.Endpoint.URL) == "" {
		return ErrEndpointRequired
	}

	parsedURL, err := url.Parse(e.Endpoint.URL)
	if err != nil {
		return ErrInvalidEndpointURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidEndpointScheme
	}

	if parsedURL.Host == "" {
		return ErrEndpointMissingHost
	}

	return nil
}

func (e ExporterConfig) validateStdoutEndpoint() error {
	// For stdout exporter, endpoint is optional (represents file path)
	if e.Endpoint.URL != "" {
		if strings.TrimSpace(e.Endpoint.URL) == "" {
			return ErrInvalidFilePath
		}

		invalidChars := []string{"\x00", "\n", "\r"}
		for _, char := range invalidChars {
			if strings.Contains(e.Endpoint.URL, char) {
				return ErrFilePathInvalidChar
			}
		}
	}

	return nil
}

func (e ExporterConfig) validateHeaders() error {
	for key, value := range e.Headers {
		if strings.TrimSpace(key) == "" {
			return ErrHeaderKeyEmpty
		}

		if strings.ContainsAny(key, " \t\n\r") {
			return ConfigError(fmt.Sprintf("header key '%s' contains invalid characters", key))
		}

		if strings.ContainsAny(value, "\n\r") {
			return ConfigError(
				fmt.Sprintf(
					"header value for key '%s' contains invalid characters",
					key,
				),
			)
		}
	}

	return nil
}

func (e ExporterConfig) Value() any {
	return e
}

// ExporterType defines the type of OTLP exporter to use.
type ExporterType uint8

func (e ExporterType) Validate() error {
	if e > ExporterTypeStdout {
		return ErrInvalidExporterType
	}

	return nil
}

func (e ExporterType) Value() any {
	return e
}

func (e ExporterType) String() string {
	switch e {
	case ExporterTypeGRPC:
		return "grpc"
	case ExporterTypeHTTP:
		return "http"
	case ExporterTypeStdout:
		return "stdout"
	default:
		return "INVALID"
	}
}

const (
	ExporterTypeGRPC ExporterType = iota
	ExporterTypeHTTP
	ExporterTypeStdout
)

const (
	serviceNameLimit   = 255
	serviceNameMinimum = 1
)

func validateServiceName(name string) error {
	name = strings.TrimSpace(name)

	// Check minimum length
	if len(name) < serviceNameMinimum {
		return ErrServiceNameEmpty
	}

	// Check maximum length (reasonable limit)
	if len(name) > serviceNameLimit {
		return ErrServiceNameTooLong
	}

	// OpenTelemetry recommends using dot notation for service names
	// Allow alphanumeric, dots, hyphens, and underscores
	for _, char := range name {
		if !isValidServiceNameChar(char) {
			return errors.Wrapf(ErrServiceNameInvalidChar, "char %s", string(char))
		}
	}

	// Service name should not start or end with dot
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return ErrServiceNameStartEndDot
	}

	// Should not have consecutive dots
	if strings.Contains(name, "..") {
		return ErrServiceNameConsecutiveDots
	}

	return nil
}

const serviceVersionLimit = 128

func isValidServiceNameChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '.' ||
		char == '-' ||
		char == '_'
}

type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

var (
	ErrNoServiceName              ConfigError = "no service name provided"
	ErrNoServiceID                ConfigError = "no service id provided"
	ErrNoServiceVersion           ConfigError = "no service version provided"
	ErrInvalidServiceName         ConfigError = "invalid service name"
	ErrInvalidServiceVersion      ConfigError = "invalid service version"
	ErrInvalidExporterConfig      ConfigError = "invalid exporter config"
	ErrServiceNameConsecutiveDots ConfigError = "service name contains consecutive dots"
	ErrServiceNameStartEndDot     ConfigError = "service name starts or ends with dot"
	ErrServiceNameInvalidChar     ConfigError = "service name contains invalid character"
	ErrServiceNameTooLong                     = ConfigError(
		"service name chars greater than " + strconv.Itoa(serviceNameLimit),
	)
	ErrServiceNameEmpty      ConfigError = "service name is empty"
	ErrServiceVersionEmpty   ConfigError = "service version is empty"
	ErrServiceVersionTooLong             = ConfigError(
		"service version chars greater than " + strconv.Itoa(serviceVersionLimit),
	)
	ErrServiceVersionInvalidChar ConfigError = "service version contains invalid character"
	ErrEndpointRequired          ConfigError = "endpoint is required for exporter"
	ErrInvalidEndpointURL        ConfigError = "invalid endpoint URL"
	ErrInvalidEndpointScheme     ConfigError = "endpoint must use http or https scheme"
	ErrEndpointMissingHost       ConfigError = "endpoint must include a host"
	ErrGRPCEndpointNoScheme      ConfigError = "gRPC endpoint must not include scheme (use host:port format)"
	ErrGRPCEndpointMissingPort   ConfigError = "gRPC endpoint must include a port (host:port format)"
	ErrInvalidFilePath           ConfigError = "if specified, endpoint must be a valid file path"
	ErrFilePathInvalidChar       ConfigError = "endpoint contains invalid character for file path"
	ErrHeaderKeyEmpty            ConfigError = "header key cannot be empty"
	ErrInvalidExporterType       ConfigError = "invalid exporter type"
)
