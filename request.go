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

import "errors"

type credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type requestItemCrypto struct {
	Pwd              string   `json:"pwd"`
	Columns          []string `json:"columns"`
	CompressionLevel int      `json:"compressionLevel,omitempty"`
}

type requestItem struct {
	Query       string                   `json:"query,omitempty"`
	Statement   string                   `json:"statement,omitempty"`
	NoFail      bool                     `json:"noFail,omitempty"`
	Values      map[string]interface{}   `json:"values,omitempty"`
	ValuesBatch []map[string]interface{} `json:"valuesBatch,omitempty"`
	Encoder     *requestItemCrypto       `json:"encoder,omitempty"`
	Decoder     *requestItemCrypto       `json:"decoder,omitempty"`
}

type request struct {
	Credentials *credentials  `json:"credentials,omitempty"`
	Transaction []requestItem `json:"transaction"`
}

type RequestBuilder struct {
	err  string
	list request
	temp *requestItem
}

type Request struct {
	req request
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{list: request{Transaction: make([]requestItem, 0)}}
}

func (rb *RequestBuilder) AddQuery(query string) *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	if rb.temp != nil {
		rb.list.Transaction = append(rb.list.Transaction, *rb.temp)
	}
	rb.temp = &requestItem{}
	rb.temp.Query = query
	return rb
}

func (rb *RequestBuilder) AddStatement(statement string) *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	if rb.temp != nil {
		rb.list.Transaction = append(rb.list.Transaction, *rb.temp)
	}
	rb.temp = &requestItem{}
	rb.temp.Statement = statement
	return rb
}

func (rb *RequestBuilder) WithNoFail() *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	rb.temp.NoFail = true
	return rb
}

func (rb *RequestBuilder) WithValues(values map[string]interface{}) *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	if values == nil {
		rb.err = "cannot specify a nil argument"
		return rb
	}
	if rb.temp.Query != "" && (rb.temp.Values != nil || rb.temp.ValuesBatch != nil) {
		rb.err = "cannot specify a batch for a query"
		return rb
	}
	if rb.temp.ValuesBatch != nil {
		rb.temp.ValuesBatch = append(rb.temp.ValuesBatch, values)
	} else if rb.temp.Values != nil {
		rb.temp.ValuesBatch = []map[string]interface{}{rb.temp.Values, values}
		rb.temp.Values = nil
	} else {
		rb.temp.Values = values
	}
	return rb
}

func (rb *RequestBuilder) WithEncoderAndCompression(password string, compressionLevel int, columns ...string) *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	if compressionLevel < 1 || compressionLevel > 19 {
		rb.err = "compressionLevel must be between 1 and 19"
		return rb
	}
	if len(columns) <= 0 {
		rb.err = "cannot specify an empty columns list"
		return rb
	}
	if rb.temp.Query != "" {
		rb.err = "cannot specify an encoder for a query"
		return rb
	}
	rb.temp.Encoder = &requestItemCrypto{
		Pwd:              password,
		CompressionLevel: compressionLevel,
		Columns:          columns,
	}
	return rb
}

func (rb *RequestBuilder) WithEncoder(password string, columns ...string) *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	if len(columns) <= 0 {
		rb.err = "cannot specify an empty columns list"
		return rb
	}
	if rb.temp.Query != "" {
		rb.err = "cannot specify an encoder for a query"
		return rb
	}
	rb.temp.Encoder = &requestItemCrypto{
		Pwd:     password,
		Columns: columns,
	}
	return rb
}

func (rb *RequestBuilder) WithDecoder(password string, columns ...string) *RequestBuilder {
	if rb.err != "" {
		return rb
	}
	if len(columns) <= 0 {
		rb.err = "cannot specify an empty columns list"
		return rb
	}
	if rb.temp.Statement != "" {
		rb.err = "cannot specify a decoder for a statement"
		return rb
	}
	rb.temp.Decoder = &requestItemCrypto{
		Pwd:     password,
		Columns: columns,
	}
	return rb
}

func (rb *RequestBuilder) Build() (*Request, error) {
	if rb.temp == nil {
		rb.err = "There are no requests"
	}
	if rb.err != "" {
		return nil, errors.New(rb.err)
	}
	rb.list.Transaction = append(rb.list.Transaction, *rb.temp)
	return &Request{rb.list}, nil
}
