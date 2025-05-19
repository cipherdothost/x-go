// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xtime

import (
	"encoding/json"
	"fmt"
	"time"

	"go.cipher.host/x/xerrors"
)

const (
	// ErrUnmarshalDuration is the error returned when a duration cannot be
	// unmarshalled.
	ErrUnmarshalDuration xerrors.Error = "could not unmarshal duration"

	// ErrParseDuration is the error returned when a duration cannot be parsed.
	ErrParseDuration xerrors.Error = "could not parse duration"
)

// Duration is a wrapper for time.Duration which supports JSON unmarshalling
// from a string.
type Duration time.Duration

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalDuration, err)
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrParseDuration, err)
	}

	*d = Duration(duration)

	return nil
}
