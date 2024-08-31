# Envelope

<p>
    <a href="https://godoc.org/github.com/eugenetriguba/envelope">
        <img src="https://godoc.org/github.com/eugenetriguba/envelope?status.svg" alt="GoDoc">
    </a>
    <a href="https://goreportcard.com/report/github.com/eugenetriguba/envelope">
        <img src="https://goreportcard.com/badge/github.com/eugenetriguba/envelope" alt="Go Report Card Badge">
    </a>
    <a href="https://codecov.io/github/eugenetriguba/envelope">
        <img src="https://codecov.io/gh/eugenetriguba/envelope/graph/badge.svg?token=1nxWoiXv66"/>
    </a>
    <img alt="Version Badge" src="https://img.shields.io/badge/version-0.1.0-blue" style="max-width:100%;">
</p>

Envelope is a Go library that enables populating Go structs from environment variables.

## Installation

To install Envelope, use the following command:

```bash
$ go get github.com/eugenetriguba/envelope
```

## Usage

Envelope supports the following types to load environment variables into for a Go struct:
- `string`
- `bool`
- `float32`, `float64`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`

Given a struct with compatible types, we can use `LoadFromEnv` to populate the struct.

`main.go`
```go
package main

import (
    "fmt"
    "os"

    "github.com/eugenetriguba/envelope"
)

type ExampleStruct struct {
    Field string `env:"EXAMPLE_FIELD"`
}

func main() {
    var example ExampleStruct
    err := envelope.LoadFromEnv(&example)
    if err != nil {
        fmt.Errorf("%w\n", err)
        os.Exit(1)
    }
    fmt.Printf("Field value: %s\n", example.Field)
}
```

Now, we can run the program to see what `Field` is set to:

```
$ go run main.go
Field value:
```

When we run it without an environment variable set, we can see it is the default Go value for that type (in this case, an empty string). However, when we set the `EXAMPLE_FIELD` environment variable, we can see the struct's `Field` field has the value set as expected:

```
$ EXAMPLE_FIELD=123 go run main.go
Field value: 123
```

