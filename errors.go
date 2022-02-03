/*
  Copyright (c) 2022-, Germano Rizzo <oss@germanorizzo.it>

  Permission to use, copy, modify, and/or distribute this software for any
  purpose with or without fee is hereby granted, provided that the above
  copyright notice and this permission notice appear in all copies.

  THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
  WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
  MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
  ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
  WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
  ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
  OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

package ws4sqlite_client

// This is an exception that wraps the error structure of ws4sqlite. See the docs at
// https://germ.gitbook.io/ws4sqlite/documentation/errors#global-errors
//
// It has fields for the error message, the index of the node that failed, and for the HTTP code.</p>
type WsError struct {
	// The index of the statement/query that failed
	QueryIndex int `json:"qryIdx"`
	// Error message
	Msg string `json:"error"`
	// HTTP code
	Code int `json:"-"`
}

func (m WsError) Error() string {
	return m.Msg
}
