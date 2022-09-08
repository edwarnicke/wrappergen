wrappergen is a library to enable generation of wrappers around binapi from Go templates.

## Quickstart

You can start a simple wrappergen project by simply adding a gen.go file to your package:
```go
//go:build tools

package main

import (
	_ "github.com/edwarnicke/wrappergen/cmd"
)

// Run using go generate -tags tools ./...
//go:generate go run github.com/edwarnicke/wrappergen/cmd"
```

and running:

```go
go generate -tags tools ./...
```

which will generate a main.go file.

Create a subdirectory templates/ and start adding *.go.tmpl files to it.  Your code generator can then be trivially consumed by downstream user with a gen.go file like:

```go
//go:build tools

package vpplink

import (
	_ "github.com/edwarnicke/vpplink/cmd"
	_ "go.fd.io/govpp/binapi"
)

// Run using go generate -tags tools ./...
//go:generate go run github.com/edwarnicke/vpplink/cmd --binapi-package "go.fd.io/govpp/binapi"
```

substituting your code generator in place of "github.com/edwarnicke/vpplink/cmd" in the above example.


