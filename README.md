# 🌱 ws4sqlite client for Go(lang)

[pkg.go.dev Docs](https://pkg.go.dev/github.com/proofrock/ws4sqlite-client-go)

This is an implementation of a client for [ws4sqlite](https://github.com/proofrock/ws4sqlite) to use with Go. It adds convenience to the communication, by not having to deal with JSON, by performing checks for the requests being well formed and by mapping errors to JDK's exceptions.

## Compatibility

Each client's minor release is guaranteed to be compatible with the matching minor release of ws4sqlite. So, for ws4sqlite's version `0.10.0`, use any of the client's `0.10.x` versions.

The library requires Go 1.17 or higher.

## Import

```bash
go get github.com/proofrock/ws4sqlite-client-go
```

# Usage

This is a translation in Go code of the "everything included" request documented in [the docs](https://germ.gitbook.io/ws4sqlite/documentation/requests). It shows the usage, overall; please refer to the [go docs]() for details.

```go
import ws4 "github.com/proofrock/ws4sqlite-client-go"

//...

// Prepare a client for the transmission. Not thread safe, but cheap to build.
cli, err := ws4.NewClientBuilder().
	WithURL("http://localhost:12321/db2").
	WithInlineAuth("myUser1", "myHotPassword").
	Build()

if err != nil {
	panic(err)
}

// Prepare the request, adding different queries/statements. See the docs for a
// detailed explanation, should be fairly 1:1 to the request at
// https://germ.gitbook.io/ws4sqlite/documentation/requests
req, err := ws4.NewRequestBuilder().
	AddQuery("SELECT * FROM TEMP").
	//
	AddQuery("SELECT * FROM TEMP WHERE ID = :id").
	WithValues(map[string]interface{}{"id": 1}).
	//
	AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (0, 'ZERO')").
	//
	AddStatement("INSERT INTO TEMP (ID, VAL) VALUES (:id, :val)").
	WithNoFail().
	WithValues(map[string]interface{}{"id": 1, "val": "a"}).
	//
	AddStatement("#Q2").
	WithValues(map[string]interface{}{"id": 2, "val": "b"}).
	WithValues(map[string]interface{}{"id": 3, "val": "c"}).
	//
	Build()

if err != nil {
	panic(err)
}

// Call ws4sqlite, obtaining a response and the status code (and a possible error)
// Status code is !=0 if the method got a response from ws4sqlite, regardless of error.
res, code, err := cli.Send(req)

// Code is 200?
if code != 200 {
	panic("There was an error, and now err can be cast to WsError")
}

if err != nil {
	wserr := err.(ws4.WsError)
	// Error possibly raised by the processing of the request.
	// It contains the same fields from
	// https://germ.gitbook.io/ws4sqlite/documentation/errors#global-errors
	fmt.Printf("HTTP Code: %d\n", wserr.Code)
	fmt.Printf("At subrequest: %d\n", wserr.RequestIdx)
	fmt.Printf("Error: %s\n", wserr.Msg) // or wserr.Error()
	panic("see above")
}

// Unpacking of the response. Every ResponseItem matches a node of the request,
// and each one has exactly one of the following fields populated:
// - Error: reason for the error, if it wasn't successful;
// - RowsUpdated: if the node was a statement and no batching was involved;
//                it's the number of updated rows;
// - RowsUpdatedBatch: if the node was a statement and a batch of values was
//                     provided; it's a slice of the numbers of updated rows
//                     for each batch item;
// - ResultSet: if the node was a query; it's a slice of maps with an item
//              per returned record, and each map has the name of the filed
//              as a key of each entry, and the value as a value.
fmt.Printf("Number of responses: %d\n", len(res.Results))

fmt.Printf("Was 1st response successful? %t\n", res.Results[0].Success)

fmt.Printf("How many records had the 1st response? %d\n", len(res.Results[0].ResultSet))

fmt.Printf("What was the first VAL returned? %s\n", res.Results[0].ResultSet[0]["VAL"])
```

The encryption extension is supported and [documented](https://pkg.go.dev/github.com/proofrock/ws4sqlite-client-go#RequestBuilder.WithDecoder). 
