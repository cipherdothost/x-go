// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xnet_test

import (
	"testing"

	"go.cipher.host/x/xnet"
)

func TestNewDialer(t *testing.T) {
	t.Parallel()

	dialer := xnet.NewDialer()

	if dialer.Timeout != xnet.DefaultDialerTimeout {
		t.Errorf("NewDialer(): expected Timeout to be %v, got %v", xnet.DefaultDialerTimeout, dialer.Timeout)
	}

	if dialer.KeepAlive != xnet.DefaultDialerKeepAlive {
		t.Errorf("NewDialer(): expected KeepAlive to be %v, got %v", xnet.DefaultDialerKeepAlive, dialer.KeepAlive)
	}
}
