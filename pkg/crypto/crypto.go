// SPDX-FileCopyrightText: 2018 Alex Edwards <https://alexedwards.net>
// SPDX-FileCopyrightText: 2024 pilcrow <pilcrowonpaper@gmail.com>
// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0
// SPDX-License-Identifier: MIT

// Package crypto implements functions related to password hash generation,
// decoding, and encoding. Some functions are modifications of the argon2id
// library by Alex Edwards and of the book "The Copenhagen Book" by Pilcrow.
package crypto

import (
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // This isn't used for storing passwords.
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"go.brokedaear.com/pkg/errors"
	"golang.org/x/crypto/argon2"
)

// ValidatePassword checks, in constant time, the equality of two
// password hashes. It accepts an input hash, which is computed from
// user input, and a stored hash, which is pre-computed, stored, and
// retrieved from a repository.
func ValidatePassword(password, storedHash string) (bool, error) {
	salt, storedHashKey, err := DecodeHash(storedHash)
	if err != nil {
		return false, err
	}
	passwordHashKey := argon2.IDKey(
		[]byte(password),
		salt,
		hashTime,
		hashMemory,
		hashThreads,
		hashKeyLen,
	)
	storedHashKeyLen := int32(len(storedHashKey)) //nolint:gosec // Integer overflow irrelevant.
	passwordKeyLen := int32(len(passwordHashKey)) //nolint:gosec // Integer overflow irrelevant.
	if subtle.ConstantTimeEq(storedHashKeyLen, passwordKeyLen) == 0 {
		return false, nil
	}
	if subtle.ConstantTimeCompare(storedHashKey, passwordHashKey) == 1 {
		return true, nil
	}
	return false, nil
}

// hashMemory is 19MB. This can be changed given hash time.
const hashMemory = 19 * 1024

// hashTime is measured in iterations.
const hashTime = 2

// hashThreads is the number of processing threads to use.
const hashThreads = 1

// hashKeyLen is the length of the generated key.
const hashKeyLen = 32

// GenerateHashedPassword generates a hash password returned as a byte slice.
// The hash is generated using Argon2.
func GenerateHashedPassword(password []byte) (string, error) {
	s, err := generateSalt()
	if err != nil {
		return "", err
	}
	salt := []byte(s)
	key := argon2.IDKey(
		password,
		salt,
		hashTime,
		hashMemory,
		hashThreads,
		hashKeyLen,
	)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)
	h := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		hashMemory,
		hashTime,
		hashThreads,
		b64Salt,
		b64Key,
	)
	return h, nil
}

// PwnHash returns a SHA-1 hex encoded hash string that can be used as input to
// haveibeenpwnd to check a password against past leaks.
//
// See: https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang
func PwnHash(password []byte) (string, error) {
	p := password
	hasher := sha1.New() //nolint:gosec // This isn't used for storing passwords.
	_, err := hasher.Write(p)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate pwn hash")
	}
	sha := hex.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

// totalBytes is the total number of bytes to generate for random values.
const totalRandHashBytes = 16

// generateSalt creates a salt that can be used in a hashing algorithm
// for password generation.
//
// See: https://thecopenhagenbook.com/random-values
func generateSalt() (string, error) {
	b := make([]byte, totalRandHashBytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate salt")
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

// totalRandStringBytes is the total number of bytes to generate for a 32 bit
// encoded random string.
const totalRandStringBytes = 15

// GenerateRandomString generates a random string of 24 alphanumeric characters.
func GenerateRandomString() (string, error) {
	customEncoding := base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ")
	bytes := make([]byte, totalRandStringBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random string")
	}
	return customEncoding.EncodeToString(bytes), nil
}

// invalidHashLength is the valid number total of keys in a stored hash.
// Stored hash values are expected to have this number of keys.
const validHashLength = 6

// DecodeHash expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
func DecodeHash(hash string) ([]byte, []byte, error) {
	var (
		err  error
		salt []byte
		key  []byte
	)
	vals := strings.Split(hash, "$")
	if len(vals) != validHashLength {
		return nil, nil, ErrInvalidHash
	}
	if vals[1] != "argon2id" {
		return nil, nil, ErrIncompatibleVariant
	}
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", hashMemory, hashTime, hashThreads)
	if err != nil {
		return nil, nil, err
	}
	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, err
	}
	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, err
	}
	return salt, key, nil
}

var (
	// ErrInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrInvalidHash = errors.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVariant is returned by ComparePasswordAndHash if the
	// provided hash was created using a unsupported variant of Argon2.
	// Currently only argon2id is supported by this package.
	ErrIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")

	// ErrIncompatibleVersion is returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")
)
