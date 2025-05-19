// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xunsafe_test

import (
	"reflect"
	"testing"

	"go.cipher.host/x/xunsafe"
)

func TestBytesToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give []byte
		want string
	}{
		{
			name: "empty byte slice",
			give: make([]byte, 0),
			want: "",
		},
		{
			name: "non-empty byte slice",
			give: []byte{72, 101, 108, 108, 111},
			want: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xunsafe.BytesToString(tt.give)
			if got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give string
		want []byte
	}{
		{
			name: "empty string",
			give: "",
			want: make([]byte, 0),
		},
		{
			name: "non-empty string",
			give: "Hello",
			want: []byte{72, 101, 108, 108, 111},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := xunsafe.StringToBytes(tt.give)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
