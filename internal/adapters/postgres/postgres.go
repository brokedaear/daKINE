// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package postgres implements a PostgreSQL database adapter.
package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.brokedaear.com/internal/common/telemetry"
	"go.brokedaear.com/internal/common/utils/loggers"
	"go.brokedaear.com/internal/core/domain"
	"go.brokedaear.com/pkg/crypto"
	"go.brokedaear.com/pkg/errors"
	"go.brokedaear.com/pkg/uuid"
)

type Postgres[T any] struct {
	db     *pgxpool.Pool
	logger loggers.Logger
	tel    telemetry.Telemetry
}

func NewPostgresDB[T any](
	ctx context.Context,
	cfg *pgxpool.Config,
	logger loggers.Logger,
	tel telemetry.Telemetry,
) (*Postgres[T], error) {
	dbpool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make new Postgres database")
	}
	return &Postgres[T]{db: dbpool, logger: logger, tel: tel}, nil
}

func (p *Postgres[T]) Close() error {
	p.logger.Info("closing postgres db connection")
	p.db.Close()
	return nil
}

// CustomerRepository implements the CustomerRepository interface for domain.Customer.
type CustomerRepository struct {
	*Postgres[domain.Customer]
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(
	ctx context.Context,
	cfg *pgxpool.Config,
	logger loggers.Logger,
	tel telemetry.Telemetry,
) (*CustomerRepository, error) {
	pg, err := NewPostgresDB[domain.Customer](ctx, cfg, logger, tel)
	if err != nil {
		return nil, err
	}
	return &CustomerRepository{Postgres: pg}, nil
}

// Insert adds a new customer to the database.
func (cr *CustomerRepository) Insert(customer *domain.Customer) error {
	ctx := context.Background()

	id, err := uuid.New()
	if err != nil {
		return errors.Wrap(err, "failed to generate customer ID")
	}
	customer.ID = id

	hashedPassword, err := crypto.GenerateHashedPassword(customer.PasswordHash)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	tx, err := cr.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		INSERT INTO users (
			id, auth0_user_id, email, email_verified, password_hash,
			total_purchases_amount, total_purchases_count, created_at,
			updated_at, last_login_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = tx.Exec(ctx, query,
		customer.ID,
		nullString(customer.AuthZeroUserID),
		customer.Email,
		customer.EmailVerified,
		hashedPassword,
		customer.TotalPurchasesAmount,
		customer.TotalPurchasesCount,
		customer.CreatedAt,
		customer.UpdatedAt,
		customer.LastLoginAt,
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert customer")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	cr.logger.Info("customer inserted successfully", "customer_id", customer.ID)
	return nil
}

// Delete soft deletes a customer by setting deleted_at timestamp.
func (cr *CustomerRepository) Delete(customer *domain.Customer) error {
	ctx := context.Background()

	tx, err := cr.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		UPDATE users 
		SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := tx.Exec(ctx, query, customer.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete customer")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("customer not found or already deleted")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	cr.logger.Info("customer deleted successfully", "customer_id", customer.ID)
	return nil
}

// UpdateInformation updates customer information (excluding password).
func (cr *CustomerRepository) UpdateInformation(customer *domain.Customer) error {
	ctx := context.Background()

	tx, err := cr.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		UPDATE users 
		SET email = $2, email_verified = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := tx.Exec(ctx, query,
		customer.ID,
		customer.Email,
		customer.EmailVerified,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update customer information")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("customer not found")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	cr.logger.Info("customer information updated successfully", "customer_id", customer.ID)
	return nil
}

// UpdatePassword updates the customer's password.
func (cr *CustomerRepository) UpdatePassword(customer *domain.Customer) error {
	ctx := context.Background()

	hashedPassword, err := crypto.GenerateHashedPassword(customer.PasswordHash)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	tx, err := cr.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		UPDATE users 
		SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := tx.Exec(ctx, query, customer.ID, hashedPassword)
	if err != nil {
		return errors.Wrap(err, "failed to update customer password")
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("customer not found")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	cr.logger.Info("customer password updated successfully", "customer_id", customer.ID)

	return nil
}

// GetByID retrieves a customer by their ID.
func (cr *CustomerRepository) GetByID(id string) (*domain.Customer, error) {
	ctx := context.Background()

	query := `
		SELECT id, auth0_user_id, email, email_verified, password_hash,
			   total_purchases_amount, total_purchases_count, created_at,
			   updated_at, last_login_at, deleted_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`

	row := cr.db.QueryRow(ctx, query, id)

	customer, err := cr.scanCustomer(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("customer not found")
		}
		return nil, errors.Wrap(err, "failed to get customer by ID")
	}

	return customer, nil
}

// GetByOAuthID retrieves a customer by their OAuth ID.
func (cr *CustomerRepository) GetByOAuthID(oauthID string) (*domain.Customer, error) {
	ctx := context.Background()

	query := `
		SELECT id, auth0_user_id, email, email_verified, password_hash,
			   total_purchases_amount, total_purchases_count, created_at,
			   updated_at, last_login_at, deleted_at
		FROM users 
		WHERE auth0_user_id = $1 AND deleted_at IS NULL`

	row := cr.db.QueryRow(ctx, query, oauthID)

	customer, err := cr.scanCustomer(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("customer not found")
		}
		return nil, errors.Wrap(err, "failed to get customer by OAuth ID")
	}

	return customer, nil
}

// GetByEmail retrieves a customer by their email address.
func (cr *CustomerRepository) GetByEmail(email string) (*domain.Customer, error) {
	ctx := context.Background()

	query := `
		SELECT id, auth0_user_id, email, email_verified, password_hash,
			   total_purchases_amount, total_purchases_count, created_at,
			   updated_at, last_login_at, deleted_at
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`

	row := cr.db.QueryRow(ctx, query, email)

	customer, err := cr.scanCustomer(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("customer not found")
		}
		return nil, errors.Wrap(err, "failed to get customer by email")
	}

	return customer, nil
}

// scanCustomer scans a database row into a domain.Customer struct.
func (cr *CustomerRepository) scanCustomer(row pgx.Row) (*domain.Customer, error) {
	var customer domain.Customer
	var auth0UserID sql.NullString
	var passwordHash string
	var deletedAt sql.NullTime

	err := row.Scan(
		&customer.ID,
		&auth0UserID,
		&customer.Email,
		&customer.EmailVerified,
		&passwordHash,
		&customer.TotalPurchasesAmount,
		&customer.TotalPurchasesCount,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.LastLoginAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if auth0UserID.Valid {
		customer.AuthZeroUserID = auth0UserID.String
	}

	if deletedAt.Valid {
		customer.DeletedAt = &deletedAt.Time
	}

	// Store the hashed password as bytes
	customer.PasswordHash = []byte(passwordHash)

	return &customer, nil
}

// nullString returns a sql.NullString for the given string value.
func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: "", Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
