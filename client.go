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

// 0.9.0-rc1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type AuthMode string

const (
	AUTH_MODE_HTTP   AuthMode = "HTTP"
	AUTH_MODE_INLINE AuthMode = "INLINE"
	AUTH_MODE_NONE   AuthMode = "NONE"
)

type Protocol string

const (
	PROTOCOL_HTTP  Protocol = "http"
	PROTOCOL_HTTPS Protocol = "https"
)

type ClientBuilder struct {
	url      string
	authMode AuthMode
	user     string
	pass     string
}

type Client struct {
	ClientBuilder
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{authMode: AUTH_MODE_NONE}
}

func (cb *ClientBuilder) WithURL(url string) *ClientBuilder {
	cb.url = url
	return cb
}

func (cb *ClientBuilder) WithURLComponents(protocol Protocol, host string, port int, databaseId string) *ClientBuilder {
	cb.url = fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, databaseId)
	return cb
}

func (cb *ClientBuilder) WithURLComponentsNoPort(protocol Protocol, host string, databaseId string) *ClientBuilder {
	cb.url = fmt.Sprintf("%s://%s/%s", protocol, host, databaseId)
	return cb
}

func (cb *ClientBuilder) WithInlineAuth(user, pass string) *ClientBuilder {
	cb.authMode = AUTH_MODE_INLINE
	cb.user = user
	cb.pass = pass
	return cb
}

func (cb *ClientBuilder) WithHTTPAuth(user, pass string) *ClientBuilder {
	cb.authMode = AUTH_MODE_HTTP
	cb.user = user
	cb.pass = pass
	return cb
}

func (cb *ClientBuilder) Build() (*Client, error) {
	if cb.url == "" {
		return nil, errors.New("no url specified")
	}
	if cb.authMode != AUTH_MODE_HTTP && cb.authMode != AUTH_MODE_NONE && cb.authMode != AUTH_MODE_INLINE {
		return nil, errors.New("invalid authMode")
	}
	if cb.authMode != AUTH_MODE_NONE && (cb.user == "" || cb.pass == "") {
		return nil, errors.New("no user or password specified")
	}
	return &Client{*cb}, nil
}

func (c *Client) Send(req *Request) (*Response, error) {
	if c.authMode == AUTH_MODE_INLINE {
		req.req.Credentials = &credentials{
			User: c.user,
			Pass: c.pass,
		}
	}

	jsonData, err := json.Marshal(req.req)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	post, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	if c.authMode == AUTH_MODE_HTTP {
		post.SetBasicAuth(c.user, c.pass)
	}
	post.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(post)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	var Res = Response{Results: make([]ResponseItem, 0)}
	for i := range res.Results {
		Ri := ResponseItem{
			Success:          res.Results[i].Success,
			Error:            res.Results[i].Error,
			RowsUpdated:      res.Results[i].RowsUpdated,
			RowsUpdatedBatch: res.Results[i].RowsUpdatedBatch,
		}
		if res.Results[i].ResultSet != nil {
			Rirs := make([]map[string]interface{}, 0)
			for i2 := range res.Results[i].ResultSet {
				Rirsi := make(map[string]interface{})
				for k, v := range res.Results[i].ResultSet[i2] {
					var v2 interface{}
					err = json.Unmarshal(v, &v2)
					if err != nil {
						return nil, err
					}
					Rirsi[k] = v2
				}
				Rirs = append(Rirs, Rirsi)
			}
			Ri.ResultSet = Rirs
		}
		Res.Results = append(Res.Results, Ri)
	}

	return &Res, nil
}
