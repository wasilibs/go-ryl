# go-ryl

go-ryl is a distribution of [ryl][1], that can be built with Go. It does not actually reimplement any
functionality of ryl in Go, instead building it into a WebAssembly binary, and
executing with the pure Go Wasm runtime [wazero][2]. This means that `go install` or `go run`
can be used to execute it, with no need to rely on separate package managers such as cargo,
on any platform that Go supports.

## Installation

Precompiled binaries are available in the [releases](https://github.com/wasilibs/go-ryl/releases).
Alternatively, install the plugin you want using `go install`.

```bash
$ go install github.com/wasilibs/go-ryl/cmd/ryl@latest
```

To avoid installation entirely, it can be convenient to use `go run`

```bash
$ go run github.com/wasilibs/go-ryl/cmd/ryl@latest .
```

[1]: https://github.com/owenlamont/ryl
[2]: https://wazero.io/
