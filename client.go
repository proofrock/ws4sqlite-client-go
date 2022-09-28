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

// 0.11.0

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Authentication mode for the database remote.
type AuthMode string

const (
	// HTTP Basic Authentication
	AUTH_MODE_HTTP AuthMode = "HTTP"
	// Credentials are inlined in the request
	AUTH_MODE_INLINE AuthMode = "INLINE"
	// No authentication
	AUTH_MODE_NONE AuthMode = "NONE"
)

// Used in URL composition
type Protocol string

const (
	// Adds http://
	PROTOCOL_HTTP Protocol = "http"
	// Adds https://
	PROTOCOL_HTTPS Protocol = "https"
)

// This class is a builder for Client instances. Once configured with the URL to
// contact and the authorization (if any), it can be used to instantiate a Client.
//
// Example:
//
//	cli, err := ws4.NewClientBuilder().
//	               WithURL("http://localhost:12321/db2").
//	               WithInlineAuth("myUser1", "myHotPassword").
//	               Build()
//
//	cli.Send(...)
type ClientBuilder struct {
	url      string
	authMode AuthMode
	user     string
	password string
}

// This struct represent a client for ws4sqlite. It can be constructed using the
// ClientBuilder struct, that configures it with the URL to contact and the authorization
// (if any). Once instantiated, it can be used to send Requests to the server.
//
// Example:
//
//	cli, err := ws4.NewClientBuilder().
//	               WithURL("http://localhost:12321/db2").
//	               WithInlineAuth("myUser1", "myHotPassword").
//	               Build()
//
//	cli.Send(...)
type Client struct {
	ClientBuilder
}

// First step when building. Generates a new ClientBuilder instance.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{authMode: AUTH_MODE_NONE}
}

// Builder methods that adds a "raw" URL for contacting the ws4sqlite remote.
func (cb *ClientBuilder) WithURL(url string) *ClientBuilder {
	cb.url = url
	return cb
}

// Builder methods that adds an URL for contacting the ws4sqlite remote, given its components.
func (cb *ClientBuilder) WithURLComponents(protocol Protocol, host string, port int, databaseId string) *ClientBuilder {
	cb.url = fmt.Sprintf("%s://%s:%d/%s", protocol, host, port, databaseId)
	return cb
}

// Builder methods that adds an URL for contacting the ws4sqlite remote, given its components but with an implicit port.
func (cb *ClientBuilder) WithURLComponentsNoPort(protocol Protocol, host string, databaseId string) *ClientBuilder {
	cb.url = fmt.Sprintf("%s://%s/%s", protocol, host, databaseId)
	return cb
}

// Builder methods that configures INLINE authentication; the remote must be configured accordingly.
func (cb *ClientBuilder) WithInlineAuth(user, password string) *ClientBuilder {
	cb.authMode = AUTH_MODE_INLINE
	cb.user = user
	cb.password = password
	return cb
}

// Builder methods that configures HTTP Basic Authentication; the remote must be configured accordingly.
func (cb *ClientBuilder) WithHTTPAuth(user, password string) *ClientBuilder {
	cb.authMode = AUTH_MODE_HTTP
	cb.user = user
	cb.password = password
	return cb
}

// Returns the Client that was built.
func (cb *ClientBuilder) Build() (*Client, error) {
	if cb.url == "" {
		return nil, errors.New("no url specified")
	}
	if cb.authMode != AUTH_MODE_HTTP && cb.authMode != AUTH_MODE_NONE && cb.authMode != AUTH_MODE_INLINE {
		return nil, errors.New("invalid authMode")
	}
	if cb.authMode != AUTH_MODE_NONE && (cb.user == "" || cb.password == "") {
		return nil, errors.New("no user or password specified")
	}
	return &Client{*cb}, nil
}

// Sends a set of requests to the remote, wrapped in a Request struct. Returns
// a matching set of responses, wrapped in a Response struct.
//
// Returns a WsError if the remote service returns a processing error. If the
// communication fails, it returns the "naked" error, so check for cast-ability.
func (c *Client) Send(req *Request) (*Response, int, error) {
	return c.SendWithContext(context.Background(), req)
}

// SendWithContext sends a set of requests to the remote with context, wrapped in a Request.
// Returns a matching set of responses, wrapped in a Response struct.
//
// Returns a WsError if the remote service returns a processing error. If the
// communication fails, it returns the "naked" error, so check for cast-ability.
func (c *Client) SendWithContext(ctx context.Context, req *Request) (*Response, int, error) {
	if c.authMode == AUTH_MODE_INLINE {
		req.req.Credentials = &credentials{
			User:     c.user,
			Password: c.password,
		}
	}

	jsonData, err := json.Marshal(req.req)
	if err != nil {
		return nil, 0, err
	}

	client := &http.Client{}
	post, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}
	if c.authMode == AUTH_MODE_HTTP {
		post.SetBasicAuth(c.user, c.password)
	}
	post.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(post)
	if err != nil {
		return nil, 0, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		wserr := WsError{}
		err = json.Unmarshal(body, &wserr)
		if err != nil {
			wserr.RequestIdx = -1
			wserr.Msg = string(body)
		}
		wserr.Code = resp.StatusCode
		return nil, resp.StatusCode, wserr
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, resp.StatusCode, err
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
						return nil, resp.StatusCode, err
					}
					Rirsi[k] = v2
				}
				Rirs = append(Rirs, Rirsi)
			}
			Ri.ResultSet = Rirs
		}
		Res.Results = append(Res.Results, Ri)
	}

	return &Res, resp.StatusCode, nil
}
