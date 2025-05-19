// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xnet

import (
	"net"
	"time"
)

const (
	// DefaultDialerTimeout is the default maximum amount of time a dial will
	// wait for a connection to complete.
	DefaultDialerTimeout = 30 * time.Second

	// DefaultDialerKeepAlive is the default amount of time a connection will be
	// kept alive.
	DefaultDialerKeepAlive = 30 * time.Second
)

// NewDialer returns a new Dialer with sane default timeout values.
func NewDialer() *net.Dialer {
	return &net.Dialer{
		Timeout:   DefaultDialerTimeout,
		KeepAlive: DefaultDialerKeepAlive,
	}
}
