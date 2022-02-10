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

package ws4sqlite_client_test

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	ws4 "github.com/proofrock/ws4sqlite-client-go"
)

func kill(cmd *exec.Cmd) {
	err := cmd.Process.Signal(syscall.SIGKILL)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	cmd := exec.Command("test/ws4sqlite-0.11.0", "--mem-db", "mydb:test/mydb.yaml", "--mem-db", "mydb2:test/mydb2.yaml")

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	defer kill(cmd)

	time.Sleep(time.Second)
	exitVal := m.Run()

	kill(cmd)

	os.Exit(exitVal)
}

func TestRequestWithHTTPAuth(t *testing.T) {
	client, err := ws4.NewClientBuilder().
		WithURLComponents(ws4.PROTOCOL_HTTP, "localhost", 12321, "mydb").
		WithHTTPAuth("myUser1", "myHotPassword").
		Build()

	if err != nil {
		t.Error(err)
	}

	request, err := ws4.NewRequestBuilder().
		AddQuery("SELECT * FROM TEMP").
		AddQuery("SELECT * FROM TEMP WHERE ID = :id ORDER BY ID ASC").
		WithValues(map[string]interface{}{"id": 1}).
		AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (0, 'ZERO')").
		AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (:id, :val)").
		WithNoFail().
		WithValues(map[string]interface{}{"id": 1, "val": "a"}).
		AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (:id, :val)").
		WithValues(map[string]interface{}{"id": 2, "val": "b"}).
		WithValues(map[string]interface{}{"id": 3, "val": "c"}).
		Build()

	if err != nil {
		t.Error(err)
	}

	res, code, err := client.Send(request)

	if err != nil {
		t.Error(err)
	}

	if code != 200 {
		t.Error("return code is not 200")
	}

	if len(res.Results) != 5 {
		t.Error("len(res.Results) != 5")
	}
	if !res.Results[0].Success {
		t.Error("!res.Results[0].Success")
	}
	if len(res.Results[0].ResultSet) != 2 {
		t.Error("len(res.Results[0].ResultSet) != 2")
	}
	if !res.Results[1].Success {
		t.Error("!res.Results[1].Success")
	}
	if len(res.Results[1].ResultSet) != 1 {
		t.Error("len(res.Results[1].ResultSet) != 1")
	}
	if res.Results[1].ResultSet[0]["VAL"] != "ONE" {
		t.Error("res.Results[1].ResultSet[0][\"VAL\"] != \"ONE\"")
	}
	if !res.Results[2].Success {
		t.Error("!res.Results[2].Success")
	}
	if *res.Results[2].RowsUpdated != 1 {
		t.Error("res.Results[2].RowsUpdated != 1")
	}
	if res.Results[3].Success {
		t.Error("res.Results[3].Success")
	}
	if res.Results[3].Error == "" {
		t.Error("res.Results[3].Error == \"\"")
	}
	if !res.Results[4].Success {
		t.Error("!res.Results[4].Success")
	}
	if len(res.Results[4].RowsUpdatedBatch) != 2 {
		t.Error("len(res.Results[4].RowsUpdatedBatch) != 2")
	}
	if res.Results[4].RowsUpdatedBatch[0] != 1 {
		t.Error("res.Results[4].RowsUpdatedBatch[0] != 1")
	}
}

func TestRequestWithInlineAuth(t *testing.T) {
	client, err := ws4.NewClientBuilder().
		WithURLComponents(ws4.PROTOCOL_HTTP, "localhost", 12321, "mydb2").
		WithInlineAuth("myUser1", "myHotPassword").
		Build()

	if err != nil {
		t.Error(err)
	}

	request, err := ws4.NewRequestBuilder().
		AddQuery("SELECT * FROM TEMP").
		AddQuery("SELECT * FROM TEMP WHERE ID = :id ORDER BY ID ASC").
		WithValues(map[string]interface{}{"id": 1}).
		AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (0, 'ZERO')").
		AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (:id, :val)").
		WithNoFail().
		WithValues(map[string]interface{}{"id": 1, "val": "a"}).
		AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (:id, :val)").
		WithValues(map[string]interface{}{"id": 2, "val": "b"}).
		WithValues(map[string]interface{}{"id": 3, "val": "c"}).
		Build()

	if err != nil {
		t.Error(err)
	}

	res, code, err := client.Send(request)

	if err != nil {
		t.Error(err)
	}

	if code != 200 {
		t.Error("return code is not 200")
	}

	if len(res.Results) != 5 {
		t.Error("len(res.Results) != 5")
	}
	if !res.Results[0].Success {
		t.Error("!res.Results[0].Success")
	}
	if len(res.Results[0].ResultSet) != 2 {
		t.Error("len(res.Results[0].ResultSet) != 2")
	}
	if !res.Results[1].Success {
		t.Error("!res.Results[1].Success")
	}
	if len(res.Results[1].ResultSet) != 1 {
		t.Error("len(res.Results[1].ResultSet) != 1")
	}
	if res.Results[1].ResultSet[0]["VAL"] != "ONE" {
		t.Error("res.Results[1].ResultSet[0][\"VAL\"] != \"ONE\"")
	}
	if !res.Results[2].Success {
		t.Error("!res.Results[2].Success")
	}
	if *res.Results[2].RowsUpdated != 1 {
		t.Error("res.Results[2].RowsUpdated != 1")
	}
	if res.Results[3].Success {
		t.Error("res.Results[3].Success")
	}
	if res.Results[3].Error == "" {
		t.Error("res.Results[3].Error == \"\"")
	}
	if !res.Results[4].Success {
		t.Error("!res.Results[4].Success")
	}
	if len(res.Results[4].RowsUpdatedBatch) != 2 {
		t.Error("len(res.Results[4].RowsUpdatedBatch) != 2")
	}
	if res.Results[4].RowsUpdatedBatch[0] != 1 {
		t.Error("res.Results[4].RowsUpdatedBatch[0] != 1")
	}
}

func TestError(t *testing.T) {
	client, err := ws4.NewClientBuilder().
		WithURLComponents(ws4.PROTOCOL_HTTP, "localhost", 12321, "mydb2").
		WithInlineAuth("myUser1", "myHotPassword").
		Build()

	if err != nil {
		t.Error(err)
	}

	request, err := ws4.NewRequestBuilder().
		AddQuery("SELENCT * FROM TEMP").
		Build()

	if err != nil {
		t.Error(err)
	}

	_, code, err := client.Send(request)

	if err == nil {
		t.Error("did not fail, but should have")
	}

	if code == 200 {
		t.Error("did return 200, but shouldn't have")
	}

	wserr, ok := err.(ws4.WsError)
	if !ok {
		t.Error("err is not a WsError")
	}

	if wserr.Code != 500 {
		t.Error("error is not 500")
	}
}
