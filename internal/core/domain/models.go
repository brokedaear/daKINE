// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package domain represents the backend domain and its types.
package domain

import (
	"time"

	"go.brokedaear.com/pkg/errors"
	"go.brokedaear.com/pkg/uuid"
)

// Customer represents a customer in the application.
type Customer struct {
	// ID is the unique UUID v7 of the customer's account in a database.
	ID string `json:"user_id"`
	// AuthZeroUserID is the Auth0 ID of the customer.
	AuthZeroUserID string `json:"-"`
	// Email is the email address of the customer.
	Email string `json:"email"`
	// EmailVerified is whether the customer has verified their email or not.
	// If the customer's email is not verified, they are not able to make
	// purchases with their account.
	EmailVerified bool `json:"email_verified"`
	// PasswordHash is the password hash of the customer's account.
	PasswordHash []byte `json:"-"`
	// TotalPurchasesAmount is the total amount a customer has spent
	// in purchases for their account's lifespan.
	TotalPurchasesAmount int `json:"total_purchases_amount"`
	// TotalPurchasesCount is the total number of purchases in a customer's
	// account.
	TotalPurchasesCount int `json:"total_purchases_count"`
	// CreatedAt is the time the customer's account was created at.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time the customer's account was updated at.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt is the time the customer's account was deleted at, which is
	// useful and necessary for auditing purposes.
	DeletedAt *time.Time `json:"-"`
	// LastLoginAt is the time the customer's account was last accessed at,
	// which is useful and necessary for auditing purposes.
	LastLoginAt time.Time `json:"last_login_at"`
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
		// TODO: might need to change LastLoginAt to when the user is actually
		// validated. Also, the user must be informed of their last login date
		// when the enter the application.
		LastLoginAt: currentTime,
	}, err
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

type ProductType struct {
	name string
}

func NewProductType(t string) (ProductType, error) {
	switch t {
	case "plugin":
		return PluginProduct, nil
	case "merchandise":
		return MerchandiseProduct, nil
	default:
		return ProductType{}, errors.New("invalid product type")
	}
}

func (p ProductType) String() string {
	return p.name
}

// Product represents one of our products, such as a plugin or a piece of
// merchandise.
type Product struct {
	// ID is the unique UUID v7 of the product in the database.
	ID string `json:"id"`
	// ProductType is the type of product the product is, such as merch
	// or plugin.
	ProductType ProductType `json:"product_type"`
	// Name is the name of the product
	Name string `json:"name"`
	// Description is the full product description.
	Description string
	// ShortDescription is a short, one line description of the product.
	ShortDescription string
	// PriceID is the PriceID of the product assigned by Stripe.
	PriceID string
	// ProductID is the Produce ID of the prodduct assigned by Stripe.
	ProductID string
	// Category ID is the ID of the category this product in the database.
	CategoryID string
	// ArtisticCredits attribute a string of credits to the product.
	ArtisticCredits string
	// CCLicense contains the Creative Commons (CC) licenses of the proudct.
	CCLicense string
	// DownloadFileName is the file name for the download.
	DownloadFileName string `json:"download_filename"`
	// CreatedAt is the date the product was created at.
	CreatedAt time.Time
	// UpdatedAt is the date the product was updated at.
	UpdatedAt time.Time
	// DeletedAt is the date the product was deleted at.
	DeletedAt time.Time
	// ReleasedAt is the product release date.
	ReleasedAt time.Time
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
		CCLicense:        "",
		DownloadFileName: "",
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
		DeletedAt:        time.Time{},
		ReleasedAt:       time.Time{},
	}
}
