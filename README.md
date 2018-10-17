# Implementation of A Password Hasher in Go

`go-password-hasher` provides an HTTP API for the encoding and persistence of passwords. `go-password-hasher`'s HTTP interface is available over port 8001.

## HTTP API

### Creating a new encoded password hash
Clients request a new encoded password hash by issuing a HTTP `POST` request to the `/hash` endpoint.

**Request:**
```http
POST http://localhost:8001/hash
password={password}
```

Hashing and encoding is performed asynchronously after the request is made. The client is returned an immediate response with a payload containing the URL where the encoded hash may be accessed and a timestamp indicating when the request may be issued.

**Response:**
```http
HTTP/1.1 202 Accepted
Content-Length: 61
Content-Type: application/json
Date: Tue, 16 Oct 2018 13:22:37 GMT

{
    "timeAvailable": "2018-10-16T07:22:42-06:00",
    "url": "/hash/1"
}

```

### Fetching an existing encoded password hash
After a client has issued a request to generate a new encoded password hash and have awaited the prescribed amount of time, they may issue a HTTP `GET` request to the URI provided in the response payload from the 
create request.

**Request:**
```http
GET http://localhost:8001/hash/{id}
```

Encoded passwords are returned to the client in plain text,

**Response**
```http
HTTP/1.1 200 OK
Content-Length: 88
Content-Type: text/plain; charset=utf-8
Date: Tue, 16 Oct 2018 13:31:33 GMT

z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg==

```

### Application Server Statistics
The `go-password-hasher` application server provides basic runtime performance statistics.

**Request:**
```http
GET http://localhost:8001/stats
```

The response payload contains two values:
* **total** - A running total of requests serviced by the application server
* **average** - The average duration (in milliseconds) for the server to respond to a request

**Response**
```http
HTTP/1.1 200 OK
Content-Length: 24
Content-Type: application/json
Date: Tue, 16 Oct 2018 13:34:45 GMT

{
    "average": 44,
    "total": 3
}
```

### Graceful, Remote Shutdown
The `go-password-hasher` application server may be remotely shutdown by issuing a HTTP `GET` request to the `/shutdown` endpoint.

```http
GET http://localhost:8001/shutdown
```

Once the request is issued, the application server will stop accepting future requests and handle all inflight requests before terminating. Once the response is received no further requests may be issued.

**Response**
```http
HTTP/1.1 200 OK
Content-Length: 0
Date: Tue, 16 Oct 2018 13:38:45 GMT
```

## Compilation
`go-password-hasher` was developed against `go 1.11`, but may compile and run on earlier versions. To compile the source code into a single, distributable binary, simply type `make build` from the project's root directory. If compilation is successful, the resulting executable can be found in the `bin` directory.

## Running

If you would rather run the application outright without AOT compilation, you can execute the following command from a shell in the project's root directory:
```shell
$ go run main.go
```
You may provide additional options and compilation flags, so please see `go help run` for details.

## Development

### Running Unit Tests
All of `go-password-hashser`'s unit tests may be found alongside their implementation files and have a `_test.go` filename suffix. Tests may be executed on a per-package basies using `go test` from the commandline:

```shell
go test github.com/mpfilbin/go-password-hasher/password
```

Alternatively, you may run all of the project's unit tests by executing `make test` from the project's root directory.