// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"slices"
	"time"

	"go.brokedaear.com/internal/core/domain"
	"go.brokedaear.com/internal/core/server"
	"go.brokedaear.com/pkg/crypto"
	"go.brokedaear.com/pkg/errors"
)

// CustomerRepository operates on data related to customer actions.
type customerRepository interface {
	Insert(*domain.Customer) error
	Delete(*domain.Customer) error
	Update(*domain.Customer) error
	GetByID(string) (*domain.Customer, error)
	GetByOAuthID(string) (*domain.Customer, error)
	GetByEmail(string) (*domain.Customer, error)
}

// CustomerService defines a service that can create, read, update, or delete
// customer data. CustomerService should be used in an authorized context.
// In other words, the data passed to the methods of CustomerService are
// assumed to be already validated. These methods are therefore already
// authenticated operations.
type CustomerService struct {
	*ServiceBase
	repo        customerRepository
	sessionRepo sessionRepository
	pwnChecker  PwnChecker[[]string]
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(
	svcBase *ServiceBase,
	repo customerRepository,
) *CustomerService {
	p := pwnCheckOnline[[]string]{
		checker: server.NewHTTPRequestClient(
			svcBase.logger,
			svcBase.tel,
			stringSliceParser[[]string]{},
		),
	}
	return &CustomerService{
		ServiceBase: svcBase, repo: repo, pwnChecker: p,
	}
}

func (c *CustomerService) Update(customer *domain.Customer) (
	*domain.Customer,
	error,
) {
	err := c.repo.Update(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// Exists returns a customer if they exist.
func (c *CustomerService) Exists(
	email string,
	auth0ID string,
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

// Delete deletes a customer by first invalidating their session and then
// removing their row in the application database.
func (c *CustomerService) Delete(customer *domain.Customer) error {
	err := c.deleteUserSession(customer)
	if err != nil {
		return err
	}
	err = c.repo.Delete(customer)
	if err != nil {
		return err
	}
	c.logger.Info("deleted customer", "customer_id", customer.ID)
	return nil
}

// SignIn signs a customer into the application.
func (c *CustomerService) SignIn(email, auth0ID, password string) (
	*domain.Customer,
	error,
) {
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
	if err != nil {
		return nil, err
	}
	if !ok {
		c.logger.Warn("incorrect user password", "customer_id", customer.ID)
		return nil, ErrCustomerLoginFailed
	}

	c.newUserSession(customer.ID)

	return customer, nil
}

// sessionDuration represents a month.
const sessionDuration = 30 * 24 * time.Hour

func (c *CustomerService) newUserSession(customerID string) {
	session := domain.NewUserSession(customerID, sessionDuration)
	c.sessionRepo.Insert(session)
}

func (c *CustomerService) deleteUserSession(customer *domain.Customer) error {
	session, ok := c.sessionRepo.GetByCustomer(customer)
	if !ok {
		return errors.New("invalid session")
	}
	return c.sessionRepo.Delete(session.Token)
}

func (c *CustomerService) validateSession(sessionID string) (bool, error) {
	session, ok := c.sessionRepo.GetByToken(sessionID)
	if !ok {
		return false, nil
	}
	if time.Now().After(session.ExpiresAt) {
		return false, errors.New("expired session")
	}
	// if time.Now().After(session.ExpiresAt.Sub(sessionExpiresIn / 2)) {
	// 	session.ExpiresAt = time.Now().Add(sessionExpiresIn)
	// }
	return true, nil
}

// SignOut signs a customer out by invalidating their login session.
func (c *CustomerService) SignOut(customer *domain.Customer) error {
	return c.deleteUserSession(customer)
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
//  2. The email is validated. The characters preceding the `@` symbol cannot
//     longer than 256 bytes.
//  3. If all is well, the new user is inserted into the repository.
//
// The auth0 sign up flow operates on the assumption that the auth0 ID
// has already been validated by some means (such as in the frontend).
// The flow follows:
//  1. The auth0 ID is checked for existence in the database. If it does exist
//     in the database and the customer making the sign-up request is who they
//     are authenticated as via Auth0, the customer data is returned--no sign up
//     takes place.
//  2. Else, the customer is inserted directly into the database.
func (c *CustomerService) SignUp(email, auth0ID, password string) (
	*domain.Customer,
	error,
) {
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
	// The way to do this is to assume that an email exists. When the user is
	// finished signing up, send a verification email to the customer. If the
	// email is legit, they will receive the email. If not, nothing happens.

	// TODO: Check if the auth0ID is alright. If an auth0 account already
	// exists and the user cookie is alright, GTFO and login.

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
// does not send a flawed sign-up request, but checking once again in
// the backend is a good sanitary habit.
func (c *CustomerService) getCustomer(email, auth0ID string) (
	*domain.Customer,
	error,
) {
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
