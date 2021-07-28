# fs-mock

A simple tool for stubbing or mocking out a JSON API server.

This is a simple, not too smart little server to help you mock sample JSON
responses from an HTTP API. This might be something you want to use when you're
developing against a client that doesn't exist yet, is too slow or expensive to
hit repeatedly, or that you want to control responses for without setting up an
integration test suite.

It isn't programmable in the sense that you prepare mocked responses in code.
Instead, you create the JSON response you want to receive for a given endpoint,
store it at a path on disk resembling the API endpoint path you want to mock,
and name the file for the HTTP verb you want to respond to.

## Installation

***From Binaries***

I'll get around to this at some point ;)

***From source***

Download this source to a computer with Go installed, navigate to it, and run:

```
go install
```

## Usage and Behavior

```
$ fs-mock -h

Usage of fs-mock:
  -host string
        host to bind to (default "0.0.0.0")
  -port int
        port to bind to (default 3000)
```

All API responses are based on the request path and method.

- Serve a file based on the API endpoint path and the request method. The file
  must be named for the method (e.g., `GET.json`, `POST.json`)
- Serve alternates by appending `?variant=$VARIANT` to API endpoint path and
  naming the file for the variant, (e.g., `GET-myvariant.json`)
- Requests for non-existing documents return a 404 with a message for the
  non-existing path requested

## Example

First, let's set up some fake data that simulates an API we want to develop
against.

```sh
$ mkdir -p ~/Desktop/fs-mock-test/
$ cd ~/Desktop/fs-mock-test
$ mkdir -p {apples/reviews,oranges/photos}
$ echo '[{"review": "pretty tasty apple", "author": "apple-enthusiast"}]' > apples/reviews/GET.json
$ echo '{"response": "OK", "reviewId": 123}' > apples/reviews/POST.json
$ echo '[{"photoUrl": "rather-large-orange.png", "uploaderId": 123}]' > oranges/photos/GET.json
$ echo '{"error": "Not Found"}' > oranges/photos/GET-error.json
```

Then, we can issue requests against those paths to simulate API responses.

First, start up `fs-mock` in that directory:

```sh
cd ~/Desktop/fs-mock-test
fs-mock -port 4000
```

Then, start issuing requests as if it were a real API with your pre-written
JSON documents serving as the mocked responses:

```sh
$ curl localhost:4000/apples/reviews
[{"review": "pretty tasty apple", "author": "apple-enthusiast"}]
$ curl localhost:4000/apples/reviews -XPOST
{"response": "OK", "reviewId": 123}
$ curl localhost:4000/oranges/photos
[{"photoUrl": "rather-large-orange.png", "uploaderId": 123}]
$ curl localhost:4000/oranges/photos?variant=error
{"error": "Not Found"}
$ curl localhost:4000/not/a/path
JSON mock not found at /home/david/Desktop/fs-mock-test/not/a/path/GET.json
```
