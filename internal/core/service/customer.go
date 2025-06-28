// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"slices"

	"go.brokedaear.com/app/domain"
	"go.brokedaear.com/internal/common/telemetry"
	"go.brokedaear.com/internal/common/utils/loggers"
	"go.brokedaear.com/internal/core/server"
	"go.brokedaear.com/pkg/crypto"
	"go.brokedaear.com/pkg/errors"
)

// CustomerRepository operates on data related to customer actions.
type CustomerRepository interface {
	CustomerAdder
	CustomerDeleter
	CustomerUpdater
	CustomerRetriever
}

type CustomerAdder interface {
	Insert(*domain.Customer) error
}

type CustomerDeleter interface {
	Delete(*domain.Customer) error
}

type CustomerUpdater interface {
	UpdateInformation(*domain.Customer) error
	UpdatePassword(*domain.Customer) error
	ResetPassword(*domain.Customer) error
}

type CustomerRetriever interface {
	GetByID(string) (*domain.Customer, error)
	GetByOAuthID(string) (*domain.Customer, error)
	GetByEmail(string) (*domain.Customer, error)
}

// CustomerService defines a service that can create, read, update, or delete
// customer data. CustomerService should be used in an authorized context.
// In other words, the data passed to the methods of CustomerService are
// assumed to be already validated. These methods are therefore already
// authenticated operations.
type CustomerService interface {
	SignIn(email, auth0ID, password string) (*domain.Customer, error)
	SignUp(email, auth0ID, password string) (*domain.Customer, error)
	Exists(email, auth0ID, password string) (*domain.Customer, error)
	Delete(customer *domain.Customer) error
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(
	repo CustomerRepository,
	tel telemetry.Telemetry,
	log loggers.Logger,
) CustomerService {
	p := pwnCheckOnline[[]string]{
		checker: server.NewHTTPRequestClient(log, tel, stringSliceParser[[]string]{}),
	}
	return &customerService{
		repo: repo, tel: tel, logger: log, pwnChecker: p,
	}
}

type customerService struct {
	repo       CustomerRepository
	tel        telemetry.Telemetry
	logger     loggers.Logger
	pwnChecker PwnChecker[[]string]
}

// Exists returns a customer if they exist.
func (c *customerService) Exists(
	email string,
	auth0ID string,
	_ string,
) (*domain.Customer, error) {
	customer, err := c.getCustomer(email, auth0ID)
	if err != nil {
		return nil, err
	}
	if customer.Email == email || customer.AuthZeroUserID == auth0ID {
		return customer, nil
	}
	return nil, ErrCustomerDoesNotExist
}

// Delete deletes a customer.
func (c *customerService) Delete(customer *domain.Customer) error {
	err := c.repo.Delete(customer)
	if err != nil {
		return err
	}
	return nil
}

// Login signs a customer into the application.
func (c *customerService) SignIn(email, auth0ID, password string) (*domain.Customer, error) {
	var (
		customer *domain.Customer
		err      error
	)
	if hasZeroValue(email, auth0ID) {
		return nil, ErrCustomerLoginFailed
	}
	if auth0ID != "" {
		customer, err = c.repo.GetByOAuthID(auth0ID)
	} else {
		customer, err = c.repo.GetByEmail(email)
	}
	if err != nil {
		return nil, ErrCustomerLoginFailed
	}
	storedHash := string(customer.PasswordHash)
	ok, err := crypto.ValidatePassword(password, storedHash)
	if err != nil || !ok {
		return nil, ErrCustomerLoginFailed
	}
	// TODO: Add the user to the active sessions database.
	return customer, nil
}

const reallyLongPasswordLength = 256

// SignUp creates a user account for a possible customer. It returns
// a the new customer and an error, if there is one.
//
// Users are able to signup in two different ways: email or auth0. Accordingly,
// SignUp has several responsibilities. The first order of business is to check
// if a customer with the credentials already exists. If the customer does not
// exist, the function chooses two paths based on whether email or an
// auth0 ID is used. An auth0 ID takes precedence over an email sign up.
//
// If email is used to sign up:
//  1. The password field is checked and validated. The password cannot be
//     longer than 256 bytes. Also, the password cannot be pwned--that means
//     it cannot exist in the "haveibeenpwned" database of leaked password
//     hashes.
//  2. The email is validated. The characters preeceeding the `@` symbol cannot
//     longer than 256 bytes.
//  3. If all is well, the new user is inserted into the repository.
//
// The auth0 sign up flow operates on the assumption that the auth0 ID
// has already been validated by some means (such as in the frontend).
// The flow follows:
//  1. The auth0 ID is checked for existence in the database. If it does exist
//     in the database and the customer making the sign up request is who they
//     are authenticated as via Auth0, the customer data is returned--no sign up
//     takes place.
//  2. Else, the customer is inserted directly into the database.
func (c *customerService) SignUp(email, auth0ID, password string) (*domain.Customer, error) {
	var (
		customer *domain.Customer
		err      error
	)

	// First, check if the user exists or not. If the user exists via email,
	// don't allow the sign up. If a user exists via Auth0, exit and start the
	// login flow.
	customer, err = c.getCustomer(email, auth0ID)
	if err != nil {
		return nil, err
	}
	if customer != nil {
		return nil, ErrCustomerAlreadyExists
	}

	// TODO: Switch based on auth0 or email. The password should NOT be checked
	// if Auth0 is used.

	// TODO: Check if the email is alright.

	// TODO: Check if the auth0ID is alright. If an auth0 account already
	// exists, GTFO and login.

	// Reject passwords greater than 256 bytes.
	if len(password) >= reallyLongPasswordLength {
		c.logger.Error("signup failed password too long")
		return nil, ErrCustomerSignUpFailed
	}

	// TODO: Validate the password. Password should be validated here.
	// Can use dropbox password validator lib.

	pwned, err := c.pwnChecker.Check(password)
	if err != nil {
		c.logger.Error("signup failed", "error", err)
		return nil, ErrCustomerSignUpFailed
	}
	if pwned {
		c.logger.Error("signup failed", "error", err)
		return nil, ErrCustomerPasswordFailed
	}

	customer, err = domain.NewCustomer(email, auth0ID, []byte(password))
	if err != nil {
		c.logger.Error("signup failed", "error", err)
		return nil, ErrCustomerSignUpFailed
	}

	err = c.repo.Insert(customer)
	if err != nil {
		c.logger.Error("signup failed", "error", err)
		return nil, ErrCustomerSignUpFailed
	}
	return customer, nil
}

// hasZeroValue checks if a slice of strings has the zero value of string types.
func hasZeroValue(vals ...string) bool {
	return slices.Contains(vals, "")
}

// getCustomer retrieves a customer from a repository based on an email or
// an auth0 ID. If the email and auth0ID are both of the type's zero value,
// an error is returned. Of course, the frontend can validate that a request
// does not send a flawed sign up request, but checking once again in
// the backend is a good sanitary habit.
func (c *customerService) getCustomer(email, auth0ID string) (*domain.Customer, error) {
	var (
		customer *domain.Customer
		err      error
	)
	if hasZeroValue(email, auth0ID) {
		return nil, ErrEmailAndAuthEmpty
	}
	if auth0ID != "" {
		customer, err = c.repo.GetByOAuthID(auth0ID)
	} else {
		customer, err = c.repo.GetByEmail(email)
	}
	return customer, err
}

var (
	ErrCustomerLoginFailed    = errors.New("customer login failed")
	ErrCustomerSignUpFailed   = errors.New("customer signup failed")
	ErrCustomerDoesNotExist   = errors.New("customer not found")
	ErrEmailAndAuthEmpty      = errors.New("email and auth0ID both zero value")
	ErrCustomerAlreadyExists  = errors.New("customer already exists")
	ErrCustomerPasswordFailed = errors.New("password found in database leak")
)
