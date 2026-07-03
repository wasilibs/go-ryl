//go:build codeanalysis

// This is a build-tag stub for the generated ryl.go (~918k lines). Under the
// `codeanalysis` tag, ryl.go is excluded from the build and this file supplies
// just the surface that host.go/runner reference, so linters and gopls can
// type-check the hand-written code without parsing the giant generated file.
// It is never compiled into a real build and its bodies are intentionally
// empty.
//
// Keep the type aliases and signatures in sync with ryl.go.
package wasm2go

import "sync"

type Module struct{}

func New(_ Xwasi_snapshot_preview1, _ Xwasi, _ Xenv) *Module { return &Module{} }

func (*Module) X_start()                       {}
func (*Module) Xmalloc(int32) int32            { return 0 }
func (*Module) X__wasilibc_cwd() *int32        { return nil }
func (*Module) Xwasi_thread_start(int32, int32) {}

type Memory = interface {
	Slice() *[]byte
	Grow(delta, max int64) int64
	Waiters() *sync.Map
}

type Xenv = interface {
	Xmemory() Memory
}

type Xwasi = interface {
	Xthread_spawn_6s4jie(v0 int32) int32
}

type Xwasi_snapshot_preview1 = interface {
	Xargs_get(v0, v1 int32) int32
	Xargs_sizes_get(v0, v1 int32) int32
	Xenviron_get(v0, v1 int32) int32
	Xenviron_sizes_get(v0, v1 int32) int32
	Xfd_close(v0 int32) int32
	Xfd_fdstat_get(v0, v1 int32) int32
	Xfd_filestat_get(v0, v1 int32) int32
	Xfd_filestat_set_size(v0 int32, v1 int64) int32
	Xfd_prestat_dir_name(v0, v1, v2 int32) int32
	Xfd_prestat_get(v0, v1 int32) int32
	Xfd_read(v0, v1, v2, v3 int32) int32
	Xfd_readdir(v0, v1, v2 int32, v3 int64, v4 int32) int32
	Xfd_seek(v0 int32, v1 int64, v2, v3 int32) int32
	Xfd_write(v0, v1, v2, v3 int32) int32
	Xpath_create_directory(v0, v1, v2 int32) int32
	Xpath_filestat_get(v0, v1, v2, v3, v4 int32) int32
	Xpath_open(v0, v1, v2, v3, v4 int32, v5, v6 int64, v7, v8 int32) int32
	Xpath_readlink(v0, v1, v2, v3, v4, v5 int32) int32
	Xpath_rename(v0, v1, v2, v3, v4, v5 int32) int32
	Xpath_unlink_file(v0, v1, v2 int32) int32
	Xproc_exit(v0 int32)
	Xrandom_get(v0, v1 int32) int32
	Xsched_yield() int32
}

func load32([]byte) uint32   { return 0 }
func store32([]byte, uint32) {}
func store64([]byte, uint64) {}
