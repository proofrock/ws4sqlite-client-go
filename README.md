# 🌱 ws4sqlite client for Go(lang)

[![javadoc](https://javadoc.io/badge2/it.germanorizzo.ws4sqlite/ws4sqlite-client-jvm/javadoc.svg)](https://javadoc.io/doc/it.germanorizzo.ws4sqlite/ws4sqlite-client-jvm)

This is an implementation of a client for [ws4sqlite](https://github.com/proofrock/ws4sqlite) to use with Go. It adds convenience to the communication, by not having to deal with JSON, by performing checks for the requests being well formed and by mapping errors to JDK's exceptions.

## Compatibility

Each client's minor release is guaranteed to be compatible with the matching minor release of ws4sqlite. So, for
ws4sqlite's version `0.9.0`, use any of the client's `0.9.x` versions.

The library requires Go 1.17 or higher.

## Import

```bash
go get github.com/proofrock/ws4sqlite-client-go
```

# Usage

This is a translation in Go code of the "everything included" request documented in 
[the docs](https://germ.gitbook.io/ws4sqlite/documentation/requests). It shows the usage, overall; please refer to the
[go docs]() for details.

```go
// Prepare a client for the transmission. It can be saved in a static final field,
// it's thread safe.
final Client cli =
        new ClientBuilder()
                .withURL("http://localhost:12321/db2")
                .withInlineAuth("myUser1", "myHotPassword")
                .build();

// Prepare the request, adding different queries/statements. See the docs for a 
// detailed explanation, should be fairly 1:1 to the request at
// https://germ.gitbook.io/ws4sqlite/documentation/requests
final Request req =
        new RequestBuilder()
                .addQuery("SELECT * FROM TEMP")

                .addQuery("SELECT * FROM TEMP WHERE ID = :id")
                .withValues(new MapBuilder().add("id", 1))

                .addStatement("INSERT INTO TEMP (ID, VAL) VALUES (0, 'ZERO')")

                .addStatement("INSERT INTO TEMP (ID, VAL) VALUES (:id, :val)")
                .withNoFail()
                .withValues(new MapBuilder().add("id", 1).add("val", "a"))

                .addStatement("#Q2")
                .withValues(new MapBuilder().add("id", 2).add("val", "b"))
                .withValues(new MapBuilder().add("id", 3).add("val", "c"))

                .build();

// Call ws4sqlite, obtaining a response
Response res;
try {
    res = cli.send(req);
} catch (ClientException ce) {
    // Exception possibly raised by the processing of the request.
    // It contains the same fields from
    // https://germ.gitbook.io/ws4sqlite/documentation/errors#global-errors
    // It is a subclass of IOException, so catch it accordingly
    System.err.format("HTTP Code: %d\n", ce.getCode());
    System.err.format("At subrequest: %d\n", ce.getQryIdx());
    System.err.format("Error: %s\n", ce.getMessage());
    return;
} catch (IOException e) {
    // This is thrown when transport errors occurs
    e.printStackTrace();
    return;
}

// Unpacking of the response. Every Response.Item matches a node of the request, 
// and each one has exactly one of the following fields populated/not null:
// - getError(): reason for the error, if it wasn't successful;
// - getRowsUpdated(): if the node was a statement and no batching was involved;
//                     it's the number of updated rows;
// - getRowsUpdatedBatch(): if the node was a statement and a batch of values was
//                          provided; it's a List of the numbers of updated rows
//                          for each batch item;
// - getResultSet(): if the node was a query; it's a List of Map()s with an item
//                   per returned record, and each map has the name of the filed
//                   as a key of each entry, and the value as a value.
System.out.format("Number of responses: %d\n",
        res.getResults().size());

System.out.format("Was 1st response successful? %s\n",
        res.getResults().get(0).isSuccess());

System.out.format("How many records had the 1st response? %d\n",
        res.getResults().get(0).getResultSet().size());

System.out.format("What was the first VAL returned? %s\n",
        res.getResults().get(0).getResultSet().get(0).get("VAL"));
```
