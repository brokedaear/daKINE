// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package domain represents the backend domain and its types.
package domain

import (
	"fmt"
	"time"

	"go.brokedaear.com/pkg/crypto"
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
	currentTime := time.Now().UTC()
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

// Order represents an entire customer order, as a result of a purchase funnel.
type Order struct {
	ID string `json:"-"`
	// UserID of the customer who made the order.
	UserID           string `json:"-"`
	StripePaymentID  string `json:"stripe_payment_id"`
	StripeCustomerID string `json:"stripe_customer_id"`
	// OrderNumber is the public-facing order number a customer can use to inquiry
	// about their order. It takes the form of:
	// BDE-YYYY-MM-DD-XXXXXXXXXXXXXXXXXXXXXXXX. "BDE" signifies our company,
	// "YYYY-MM-DD" signifies the order date, and "X..." signifies the order
	// number.
	OrderNumber string `json:"order_number"`
	// Items are the items in the order.
	Items []LineItem `json:"items"`
	// GrandTotal is the total amount to be paid.
	GrandTotal int `json:"grand_total"`
	// CurrencyID is the three character currency identifier, like USD or CNY.
	CurrencyID string            `json:"currency_id"`
	Status     FulfillmentStatus `json:"status"`
	// CreatedAt is the date the order was created at.
	CreatedAt time.Time
	// UpdatedAt is the date the order was updated at.
	UpdatedAt time.Time
	// CompletedAt is the date the order was completed at; the fulfillment date.
	// For plugins, this is near immediately. For merchandise, this is when the
	// product delivery occurs.
	CompletedAt *time.Time
	// DeletedAt is the date the order was deleted at.
	DeletedAt *time.Time
}

// NewOrder creates a new customer order.
func NewOrder(items ...LineItem) (*Order, error) {
	now, id, err := newTimeWithID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make new order")
	}
	orderID, err := NewOrderNumber(*now)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make new order")
	}
	return &Order{
		ID:          id,
		Items:       items,
		OrderNumber: orderID,
		Status:      PendingStatus,
		CreatedAt:   *now,
		UpdatedAt:   *now,
		CompletedAt: nil,
		DeletedAt:   nil,
	}, nil
}

// LineItem is a product with order related details.
type LineItem struct {
	// ID is the UUID v7 internal product ID.
	ID string `json:"-"`

	// Product is the product on the line.
	Product Product `json:"product"`

	// Quantity is the total number of the product the customer intends to purchase.
	Quantity int `json:"quantity"`

	// Status is the status of the line item. Line items can have separate
	// statuses, because a customer can order a plugin and merchandise in the same
	// order. A plugin can be delivered immediately, while a piece of merchandise
	// will obviously take time to be delivered.
	Status FulfillmentStatus `json:"status"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

func NewLineItem(product Product, quantity int) (*LineItem, error) {
	now, id, err := newTimeWithID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make new order")
	}
	return &LineItem{
		ID:        id,
		Product:   product,
		Quantity:  quantity,
		Status:    PendingStatus,
		CreatedAt: *now,
		UpdatedAt: *now,
		DeletedAt: nil,
	}, nil
}

func NewOrderNumber(t time.Time) (string, error) {
	const prefix = "BDE"
	date := t.Format("2006-01-02")
	id, err := crypto.GenerateRandomString()
	if err != nil {
		return "", errors.Wrap(err, "failed to generate new order number")
	}
	return fmt.Sprintf("%s-%s-%s", prefix, date, id), nil
}

var (
	//nolint:gochecknoglobals // These simulate enums.
	PendingStatus = FulfillmentStatus{status: "pending"}
	//nolint:gochecknoglobals // These simulate enums.
	ProcessingStats = FulfillmentStatus{status: "processing"}
	//nolint:gochecknoglobals // These simulate enums.
	CompletedStatus = FulfillmentStatus{status: "completed"}
	//nolint:gochecknoglobals // These simulate enums.
	FailedStatus = FulfillmentStatus{status: "failed"}
	//nolint:gochecknoglobals // These simulate enums.
	RefundedStatus = FulfillmentStatus{status: "refunded"}
	//nolint:gochecknoglobals // These simulate enums.
	CancelledStatus = FulfillmentStatus{status: "cancelled"}
)

// FulfillmentStatus represents the status of a line item or order.
type FulfillmentStatus struct {
	status string
}

func (f FulfillmentStatus) String() string {
	return f.status
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
	//nolint:gochecknoglobals // These simulate enums.
	InvalidProduct = ProductType{
		name: "",
	}
)

// ProductType is a pseudo-enum that is the type of product sold. There are only
// two possible values: "plugin" and "merchandise".
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
		return InvalidProduct, errors.New("invalid product type")
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
	// Price is the price of the product.
	Price int
	// ProductID is the Produce ID of the prodduct assigned by Stripe.
	ProductID string
	// Category ID is the ID of the category this product in the database.
	CategoryID string
	// Credits contains the authors of assert licences and the license, in the
	// format of `<author_name>-<SPDX_identifier>;`
	Credits string
	// DownloadFileName is the file name for the download.
	DownloadFileName string `json:"download_filename"`
	// DownloadChecksum is the checksum of the product if it is a file.
	DownloadChecksum string
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
		Price:            0,
		ProductID:        "",
		CategoryID:       "",
		Credits:          "",
		DownloadFileName: "",
		DownloadChecksum: "",
		CreatedAt:        time.Time{},
		UpdatedAt:        time.Time{},
		DeletedAt:        time.Time{},
		ReleasedAt:       time.Time{},
	}
}

func newTimeWithID() (*time.Time, string, error) {
	currentTime := time.Now().UTC()
	internalID, err := uuid.New()
	if err != nil {
		return nil, "", err
	}
	return &currentTime, internalID, nil
}
