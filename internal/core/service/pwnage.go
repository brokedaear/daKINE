// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"strings"

	"go.brokedaear.com/internal/core/server"
	"go.brokedaear.com/pkg/collections"
	"go.brokedaear.com/pkg/crypto"
)

// PwnChecker checks if a password has been pwned.
type PwnChecker[T any] interface {
	Check(password string) (bool, error)
}

type pwnCheckOnline[T any] struct {
	checker server.HTTPClient[[]string]
}

func (p pwnCheckOnline[T]) Check(password string) (bool, error) {
	pwnHash, err := crypto.PwnHash([]byte(password))
	if err != nil {
		return false, err
	}
	target := pwnHash[:5]
	pwnURL := "https://api.pwnedpasswords.com/range/" + target
	// Add request padding in the response. This adds 80% more bytes to the
	// response, So I don't think this is necessarily needed. However, I
	// don't think our scale will be monstrous, so this may be fine in terms
	// of ingress/egress stuff. I think its fine...
	//
	// See:
	// https://www.troyhunt.com/enhancing-pwned-passwords-privacy-with-padding/
	// TODO: Context.
	headers := make(map[string]string)
	headers["Add-Padding"] = "true"
	leakedHashes, err := p.checker.Get(context.TODO(), pwnURL, headers)
	if err != nil {
		return false, err
	}
	leakedHashes = collections.Filter(leakedHashes, func(target, item string) bool {
		return target == item[:5]
	}, target)
	_, ok := collections.Find(leakedHashes, func(a1, a2 string) bool {
		return a1 == a2
	}, pwnHash)
	return ok, nil
}

type stringSliceParser[T []string] struct{}

// Parse reads a body of []byte into a string slice.
func (s stringSliceParser[T]) Parse(body []byte) (T, error) {
	v := strings.Split(string(body), "\n")
	return v, nil
}
