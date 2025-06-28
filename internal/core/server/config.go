// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net"
	"path"
	"strconv"
	"strings"

	"go.brokedaear.com/internal/common/telemetry"
	"go.brokedaear.com/pkg/errors"
	"go.brokedaear.com/pkg/validator"
)

// config defines a default server configuration.
type config struct {
	// addr is the Address on which to bind the application.
	addr Address

	// port number to bind to for the application.
	port Port

	// socketPath is the path of a Unix socket that a server may use to communicate.
	socketPath SocketPath

	// version is the version of the software.
	version Version

	telemetry telemetry.Telemetry
}

func newConfig(opts ...ConfigOpts) (*config, error) {
	var err error
	config := newDefaultConfig()
	for _, opt := range opts {
		err = opt(config)
	}
	if err != nil {
		return nil, err
	}
	err = config.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make config")
	}
	return config, nil
}

const (
	defaultAddress = "127.0.0.1"
	defaultPort    = 8080
	defaultVersion = "0.0.1"
)

func newDefaultConfig() *config {
	return &config{
		addr:       defaultAddress,
		port:       Port(defaultPort),
		socketPath: "",
		version:    defaultVersion,
		telemetry:  nil,
	}
}

func (c config) Validate() error {
	err := validator.Check(
		c.addr,
		c.port,
		c.socketPath,
		c.version,
	)
	if err != nil {
		return errors.Wrap(err, "config failed validation")
	}
	return nil
}

func (c config) Value() any {
	return c
}

// newURIAddress creates a bindable URI address from Addr and Port.
func (c config) newURIAddress() (string, error) {
	err := validator.Check(c.addr, c.port)
	if err != nil {
		return "", errors.Wrap(err, "failed to create new URI address from config")
	}
	address := net.JoinHostPort(c.addr.String(), c.port.String())
	return address, err
}

type ConfigOpts func(*config) error

func WithAddress(addr string) ConfigOpts {
	return func(c *config) error {
		c.addr = Address(addr)
		return nil
	}
}

func WithPort(port uint16) ConfigOpts {
	return func(c *config) error {
		c.port = Port(port)
		return nil
	}
}

func WithSocketPath(path string) ConfigOpts {
	return func(c *config) error {
		c.socketPath = SocketPath(path)
		return nil
	}
}

func WithVersion(vers string) ConfigOpts {
	return func(c *config) error {
		c.version = Version(vers)
		return nil
	}
}

func WithTelemetry(t telemetry.Telemetry) ConfigOpts {
	return func(c *config) error {
		c.telemetry = t
		return nil
	}
}

// Address represents a layer 4 OSI Address. An address must only be either an
// IP address, a domain name followed by a TLD, or a path to a Unix socket.
type Address string

func (a Address) String() string {
	return string(a)
}

func (a Address) Validate() error {
	const (
		colon        = ":"
		space        = " "
		forwardSlash = "/"
	)

	addr := a.String()

	if len(addr) == 0 {
		return ErrInvalidAddressLength
	}

	if strings.Contains(addr, colon) {
		return ErrInvalidAddressColon
	}

	if strings.Contains(addr, space) {
		return ErrInvalidAddressSpace
	}

	if strings.Contains(addr, forwardSlash) {
		return ErrInvalidAddressWithPath
	}

	return nil
}

func (a Address) Value() any {
	return a.String()
}

// Port represents a layer 4 OSI Port. A port with a zero value
// will be ignored. Abstain from assigning a value to a port when using
// a Unix socket to communicate data.
type Port uint16

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

func (p Port) Validate() error {
	if (p > 0 && p < 1024) || p >= 65534 {
		return errors.Errorf("invalid port %d must be [1024, 65535)", p)
	}

	return nil
}

func (p Port) Value() any {
	return uint16(p)
}

type SocketPath string

func (s SocketPath) String() string {
	return string(s)
}

func (s SocketPath) Validate() error {
	p := path.Ext(s.String())
	if p != "sock" && len(p) > 0 {
		return errors.New("invalid socket")
	}
	return nil
}

func (s SocketPath) Value() any {
	return s
}

type Version string

func (v Version) String() string {
	return string(v)
}

func (v Version) Validate() error {
	const expectedVersionParts = 3 // major.minor.patch
	elements := strings.Split(v.String(), ".")
	if len(elements) != expectedVersionParts {
		return ErrInvalidVersionFormat
	}

	for _, element := range elements {
		n, err := strconv.Atoi(element)
		if err != nil {
			return ErrInvalidVersionChars
		}

		if n < 0 {
			return ErrInvalidVersionSign
		}
	}

	return nil
}

func (v Version) Value() any {
	return v.String()
}

type ConfigError string

func (c ConfigError) Error() string {
	return string(c)
}

const (
	ErrInvalidPortRange       ConfigError = "port range must be [1024, 65535)"
	ErrInvalidVersionFormat   ConfigError = "version must be of the format x.x.x"
	ErrInvalidVersionChars    ConfigError = "version must only be an unsigned integer"
	ErrInvalidVersionSign     ConfigError = "version must be >= 0"
	ErrInvalidAddressLength   ConfigError = "address length must be greater than 0"
	ErrInvalidAddressColon    ConfigError = "address must not contain a colon"
	ErrInvalidAddressSpace    ConfigError = "address must not contain a space"
	ErrInvalidAddressWithPath ConfigError = "address must not contain a path"
	ErrInvalidVersionAlpha    ConfigError = "version cannot contain alpha chars"
	ErrInvalidEmptySocketPath ConfigError = "socket path cannot be empty"
)
