package wasm

import _ "embed"

//go:embed memory.wasm
var Memory []byte

//go:embed ryl.wasm
var Ryl []byte
