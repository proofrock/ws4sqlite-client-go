/*
  Copyright (c) 2022-, Germano Rizzo <oss /AT/ germanorizzo /DOT/ it>

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

import "encoding/json"

type responseItem struct {
	Success          bool                         `json:"success"`
	RowsUpdated      *int64                       `json:"rowsUpdated"`
	RowsUpdatedBatch []int64                      `json:"rowsUpdatedBatch"`
	ResultSet        []map[string]json.RawMessage `json:"resultSet"`
	Error            string                       `json:"error"`
}

type response struct {
	Results []responseItem `json:"results"`
}

// Single response coming from ws4sqlite. Every ResponseItem matches a node of the request,
// and each one has exactly one of the following fields populated/not null:
//
// - Error: reason for the error, if it wasn't successful;
//
// - RowsUpdated: if the node was a statement and no batching was involved; it's the number
// of updated rows;
//
// - RowsUpdatedBatch: if the node was a statement and a batch of values was provided; it's
// a slice of the numbers of updated row for each batch item;
//
// - ResultSet: if the node was a query; it's a slice of maps with an item per returned
// record, and each map has the name of the filed as a key of each entry, and the value as a value.
type ResponseItem struct {
	// Was the request successful?
	Success bool
	// If the node was a statement and no batching was involved, it's the number of updated
	// rows
	RowsUpdated *int64
	// If the node was a statement and a batch of values was provided, it's a slice of the
	// numbers of updated rows for each batch item
	RowsUpdatedBatch []int64
	// If the node was a query, it's a slice of maps with an item per returned record, and
	// each map has the name of the filed as a key of each entry, and the value as a value
	ResultSet []map[string]interface{}
	// Reason for the error, if the request wasn't successful
	Error string
}

// Response coming from the endpoint, that is a list of single responses
// matching the list of request that were submitted. The single responses
// are of type ResponseItem.
type Response struct {
	// Slice with the results, each one is a ResponseItem
	Results []ResponseItem
}
