// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

// Package xargon2 provides functions and utilities to extend [Go's argon2
// module].
//
// [Go's argon2 module]: https://pkg.go.dev/golang.org/x/crypto/argon2
package xargon2

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"go.cipher.host/x/xcrypto"
	"go.cipher.host/x/xcrypto/xrand"
	"go.cipher.host/x/xerrors"
	"golang.org/x/crypto/argon2"
)

const (
	// ErrInvalidPassword is returned when the provided password does not match
	// the stored hash during verification.
	ErrInvalidPassword xerrors.Error = "invalid password"

	// ErrInvalidHash is returned when the provided hash string is not in the
	// expected format.
	ErrInvalidHash xerrors.Error = "invalid hash format"

	// ErrParseHashParameters is returned when CompareHashAndPassword fails to
	// parse the parameters section of a hash string.
	ErrParseHashParameters xerrors.Error = "failed to parse hash parameters"

	// ErrDecode is returned when CompareHashAndPassword fails to decode the
	// base64-encoded hash or salt.
	ErrDecode xerrors.Error = "failed to decode hash or salt"
)

// DefaultSaltAndPepperLength is the default length in bytes for salt and pepper
// values used by GenerateFromPassword when no salt is provided.
const DefaultSaltAndPepperLength int = 16

// Parameters holds the configuration parameters for the Argon2 hashing
// algorithm.
type Parameters struct {
	// Threads is the number of threads (parallelism parameter) to use. It
	// should typically be set to the number of available CPU cores.
	Threads uint8

	// Time is the computational cost factor (number of iterations). Higher
	// values increase security at the cost of longer computation time.
	Time uint32

	// Memory is the memory usage in KiB. Higher values increase security but
	// require more memory.
	Memory uint32

	// KeyLen is the length of the generated key (hash) in bytes.
	KeyLen uint32
}

// RecommendedParameters returns a Parameters instance with [values recommended
// by RFC9106] for password hashing. These parameters are suitable for
// general-purpose password hashing on modern hardware.
//
// The returned parameters use 4 threads, time cost of 1, memory cost of 2GiB,
// and a key length of 32 bytes.
//
// [values recommended by RFC9106]: https://datatracker.ietf.org/doc/html/rfc9106#section-7.4
func RecommendedParameters() Parameters {
	return Parameters{
		Threads: 4,
		Time:    1,
		Memory:  2 * 1024 * 1024,
		KeyLen:  32,
	}
}

// MemoryConstrainedParameters returns a Parameters instance suitable for
// memory-constrained environments [as recommended by RFC9106].
//
// The returned parameters use 4 threads, time cost of 3, memory cost of 64MiB,
// and a key length of 32 bytes. The increased time cost compensates for the
// reduced memory usage.
//
// [as recommended by RFC9106]: https://datatracker.ietf.org/doc/html/rfc9106#section-7.4
func MemoryConstrainedParameters() Parameters {
	return Parameters{
		Threads: 4,
		Time:    3,
		Memory:  64 * 1024,
		KeyLen:  32,
	}
}

// CompareHashAndPassword compares an argon2id hashed password with its possible
// plaintext equivalent.
//
// The provided hash must be in the standard PHC string format:
// $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
//
// If a pepper was used during hash generation, the same pepper must be provided
// for successful comparison.
func CompareHashAndPassword(hashedPassword, password string, pepper []byte) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return ErrInvalidHash
	}

	var (
		memory, time uint32
		threads      uint8
	)

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrParseHashParameters, err)
	}

	var (
		saltBase64 = parts[4]
		hashBase64 = parts[5]
	)

	salt, err := base64.RawStdEncoding.DecodeString(saltBase64)
	if err != nil {
		return fmt.Errorf("%w: salt: %w", ErrDecode, err)
	}

	hashBytes, err := base64.RawStdEncoding.DecodeString(hashBase64)
	if err != nil {
		return fmt.Errorf("%w: hash: %w", ErrDecode, err)
	}

	var (
		passwordBytes    = []byte(password)
		combinedPassword = passwordBytes
	)

	if pepper != nil {
		combinedPassword = append(combinedPassword, pepper...)
	}

	newHash := argon2.IDKey(combinedPassword, salt, time, memory, threads, uint32(len(hashBytes))) //nolint:gosec // seems like a safe conversion since Argon2 hashes are constrained to small sizes, making uint32 overflow impractical in practice

	if subtle.ConstantTimeCompare(hashBytes, newHash) != 1 {
		return ErrInvalidPassword
	}

	return nil
}

// GenerateFromPassword returns a secure Argon2id hash of the given password
// using the specified parameters.
//
// The hash is returned in the PHC string format:
// $argon2id$v=19$m=<memory>,t=<iterations>,p=<parallelism>$<salt>$<hash>
//
// If salt is nil, a cryptographically secure random salt of
// DefaultSaltAndPepperLength bytes is generated. For most applications, nil
// should be passed to allow the function to generate a secure random salt.
//
// The optional pepper is a secret key that can be appended to the password
// before hashing for additional security. Unlike the salt, the pepper is not
// stored in the resulting hash string and must be provided separately during
// verification.
func GenerateFromPassword(password string, parameters Parameters, salt, pepper []byte) string {
	if salt == nil {
		salt = xrand.Bytes(DefaultSaltAndPepperLength)
	}

	defer xcrypto.ZeroMemory(salt)

	var (
		passwordBytes    = []byte(password)
		combinedPassword = passwordBytes
	)

	if pepper != nil {
		combinedPassword = append(combinedPassword, pepper...)
	}

	defer func() {
		xcrypto.ZeroMemory(passwordBytes)
		xcrypto.ZeroMemory(combinedPassword)
	}()

	hash := argon2.IDKey(
		combinedPassword,
		salt,
		parameters.Time,
		parameters.Memory,
		parameters.Threads,
		parameters.KeyLen,
	)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		parameters.Memory, parameters.Time, parameters.Threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}
