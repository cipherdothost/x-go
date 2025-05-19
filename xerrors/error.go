// SPDX-FileCopyrightText: 2025 The Cipher Host Team <team@cipher.host>
//
// SPDX-License-Identifier: MIT

package xerrors

// Error is an [immutable error] type.
//
// [immutable error]: https://dave.cheney.net/2016/04/07/constant-errors
type Error string

// Error implements the error interface for Error.
func (e Error) Error() string {
	return string(e)
}
