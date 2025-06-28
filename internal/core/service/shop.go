// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"go.brokedaear.com/app/domain"
)

type WebshopCustomerHandler interface {
	CustomerUpdater
	CustomerRetriever
}

// PaymentProcessor handles and processes payments.
type PaymentProcessor interface {
	Refund()
	Pay()
}

// WebshopService enables customers to purchase and refund their
// products.
type WebshopService interface {
	Purchase(product *domain.Product) error
}

func NewWebshopService(processor PaymentProcessor, repo WebshopCustomerHandler) WebshopService {
	return &webshopService{
		paymentProcessor: processor,
		repo:             repo,
	}
}

type webshopService struct {
	paymentProcessor PaymentProcessor
	repo             WebshopCustomerHandler
}

func (w webshopService) Purchase(_ *domain.Product) error {
	return nil
}
