# Https
A simple to use Golang http request library

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/mustache-cn/https) [![GoDoc](https://pkg.go.dev/badge/github.com/mustache-cn/https?utm_source=godoc)](https://godoc.org/github.com/mustache-cn/https)[![License MIT](https://img.shields.io/github/license/mustache-cn/https)](https://github.com/mustache-cn/https)


License
======

https is licensed under the MIT License, Version 2.0. See [LICENSE](LICENSE) for the full license text

Features
========

- Chain call, easy to assemble request
- Responses can be serialized into JSON
- Easy file downloads
- Support for the following HTTP verbs `GET, HEAD, POST, PUT, DELETE, PATCH, OPTIONS`

Install
=======
`go get -u github.com/mustache-cn/https`

Usage
======
`import "github.com/mustache-cn/https"`

Basic Examples
=========
Basic GET request:

```go
	response, err := https.NewClient(url).
		AddParam("key", "value").
		AddHeader("key","value").
		Get()
	// request failed
	if err != nil {
		return result, err
	}
```

If an error occurs all of the other properties and methods of a `Response` will be `nil`

Quirks
=======
## Request Quirks

When passing parameters to be added to a URL, if the URL has existing parameters that *_contradict_* with what has been passed within `Params` – `Params` will be the "source of authority" and overwrite the contradicting URL parameter.

Lets see how it works...

```go
https.NewClient(url).AddParam("key", "value")
```

## Response Quirks

Order matters! This is because `https.Response` is implemented as an `io.ReadCloser` which proxies the *https.Response.Body* `io.ReadCloser` interface. It also includes an internal buffer for use in `Response.String()` and `Response.Bytes()`.

Here are a list of methods that consume the *http.Response.Body* `io.ReadCloser` interface.

- Response.JSON
- Response.DownloadToFile

The following methods make use of an internal byte buffer

- Response.String
- Response.Bytes

In the code below, once the file is downloaded – the `Response` struct no longer has access to the request bytes

```go
	response, err := https.NewClient(url).
		AddParam("key", "value").
		AddHeader("key","value").
		Get()
	// request failed
	if err != nil {
		return result, err
	}

// At this point the .String and .Bytes method will return empty responses

response.Bytes() == nil // true
response.String() == "" // true

```

But if we were to call `response.Bytes()` or `response.String()` first, every operation will succeed until the internal buffer is cleared:

```go
	response, err := https.NewClient(url).
		AddParam("key", "value").
		AddHeader("key","value").
		Get()
	// request failed
	if err != nil {
		return result, err
	}

// This call to .Bytes caches the request bytes in an internal byte buffer – which can be used again and again until it is cleared
response.Bytes() == `file-bytes`
response.String() == "file-string"

// This will work because it will use the internal byte buffer
if err := response.DownloadToFile("randomFile"); err != nil {
	log.Println("Unable to download file: ", err)
}

// Now if we clear the internal buffer....
response.ClearInternalBuffer()

// At this point the .String and .Bytes method will return empty responses

response.Bytes() == nil // true
response.String() == "" // true
```