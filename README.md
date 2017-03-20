# Grove client

Go client library for interacting with [Grove](https://github.com/t11e/grove).

# Usage

TODO

# Contributions

Clone this repository into your GOPATH (`$GOPATH/src/github.com/t11e/`)
and use [Glide](https://github.com/Masterminds/glide) to install its dependencies.

```sh
brew install glide
go get github.com/t11e/go-groveclient
cd "$GOPATH"/src/github.com/t11e/go-groveclient
glide install --strip-vendor
```

You can then run the tests:

```sh
go test $(go list ./... | grep -v /vendor/)
```

There is no need to use `go install` as any project that requires this library
can include it as a dependency like so:

```sh
cd my_other_project
glide get --strip-vendor github.com/t11e/go-groveclient
```

If you change any of the interfaces that have a mock in `mocks/` directory be sure to execute
`go generate` and check in the updated mock files.
