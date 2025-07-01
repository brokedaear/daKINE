// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"errors"
	"time"

	"go.brokedaear.com/internal/core/domain"
)

type sessionRepository interface {
	GetByToken(token string) (session *domain.UserSession, ok bool)
	GetByCustomer(customer *domain.Customer) (session *domain.UserSession, ok bool)
	Insert(customerSession *domain.UserSession)
	Update()
	Delete(token string) error
}

// SessionService manages a customer account session.
type SessionService struct {
	*ServiceBase
	repo sessionRepository
}

func NewSessionService(svcBase *ServiceBase, repo sessionRepository) *SessionService {
	return &SessionService{
		ServiceBase: svcBase,
		repo:        repo,
	}
}

func (s *SessionService) NewSession(userID string) *domain.UserSession {
	session := domain.NewUserSession(userID, sessionDuration)
	s.repo.Insert(session)
	return session
}

// Validate checks if a session is valid, given its token.
func (s *SessionService) Validate(token string) (bool, error) {
	session, ok := s.repo.GetByToken(token)
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
