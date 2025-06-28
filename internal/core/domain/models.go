// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package domain represents the backend domain and its types.
package domain

import (
	"time"

	"go.brokedaear.com/pkg/uuid"
)

// Customer represents a customer in the application.
type Customer struct {
	// ID is the unique UUID v7 of the customer's account in a database.
	ID string
	// AuthZeroUserID is the Auth0 ID of the customer.
	AuthZeroUserID string
	// Email is the email address of the customer.
	Email string `json:"email"`
	// EmailVerified is whether the customer has verified their email or not.
	// If the customer's email is not verified, they are not able to make
	// purchases with their account.
	EmailVerified bool
	// PasswordHash is the password hash of the customer's account.
	PasswordHash []byte `json:"-"`
	// TotalPurchasesAmount is the total amount a customer has spent
	// in purchases for their account's lifespan.
	TotalPurchasesAmount int
	// TotalPurchasesCount is the total number of purchases in a customer's
	// account.
	TotalPurchasesCount int
	// CreatedAt is the time the customer's account was created at.
	CreatedAt time.Time
	// UpdatedAt is the time the customer's account was updated at.
	UpdatedAt time.Time
	// DeletedAt is the time the customer's account was deleted at, which is
	// useful and necessary for auditing purposes.
	DeletedAt *time.Time
	// LastLoginAt is the time the customer's account was last accessed at,
	// which is useful and necessary for auditing purposes.
	LastLoginAt time.Time
}

func NewCustomer(email, auth0ID string, passwordHash []byte) (*Customer, error) {
	currentTime := time.Now()
	id, err := uuid.New()
	if err != nil {
		return nil, err
	}
	return &Customer{
		ID:                   id,
		AuthZeroUserID:       auth0ID,
		Email:                email,
		EmailVerified:        false,
		PasswordHash:         passwordHash,
		TotalPurchasesAmount: 0,
		TotalPurchasesCount:  0,
		CreatedAt:            currentTime,
		UpdatedAt:            currentTime,
		DeletedAt:            nil,
		LastLoginAt:          currentTime, // TODO: might need to change this to when the user is actually validated.
	}, err
}

type ProductType struct {
	name string
}

func (p ProductType) String() string {
	return p.name
}

var (
	//nolint:gochecknoglobals // These simulate enums.
	MerchandiseProduct = ProductType{
		name: "merchandise",
	}
	//nolint:gochecknoglobals // These simulate enums.
	PluginProduct = ProductType{
		name: "plugin",
	}
)

type Product struct {
	ID               string
	ProductType      ProductType
	Name             string
	Description      string
	ShortDescription string
	PriceID          string
	ProductID        string
	CategoryID       string
	ArtisticCredits  string
	TechnicalCredits string
	CCLicense        string
	DownloadFileName string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
	ReleasedAt       time.Time
}

func NewProduct(prodType ProductType, prodID, prodName string) *Product {
	return &Product{
		ID:               prodID,
		ProductType:      prodType,
		Name:             prodName,
		Description:      "",
		ShortDescription: "",
		PriceID:          "",
		ProductID:        "",
		CategoryID:       "",
		ArtisticCredits:  "",
		TechnicalCredits: "",
		CCLicense:        "",
		DownloadFileName: "",
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
		DeletedAt:        time.Time{},
		ReleasedAt:       time.Time{},
	}
}
