package runner

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
	"github.com/wasilibs/wazero-helpers/allocator"

	"github.com/wasilibs/go-ryl/internal/wasm"
)

func Run(name string, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer, cwd string) int {
	ctx := context.Background()

	ctx = experimental.WithMemoryAllocator(ctx, allocator.NewNonMoving())

	rtCfg := wazero.NewRuntimeConfig().
		WithCoreFeatures(api.CoreFeaturesV2 | experimental.CoreFeaturesThreads).
		WithMemoryLimitPages(uint32(65536))
	uc, err := os.UserCacheDir()
	if err == nil {
		cache, err := wazero.NewCompilationCacheWithDir(filepath.Join(uc, "com.github.wasilibs"))
		if err == nil {
			rtCfg = rtCfg.WithCompilationCache(cache)
		}
	}
	rt := wazero.NewRuntimeWithConfig(ctx, rtCfg)

	code, err := rt.CompileModule(ctx, wasm.Ryl)
	if err != nil {
		panic(err)
	}

	wasi_snapshot_preview1.MustInstantiate(ctx, rt)

	if _, err := rt.InstantiateWithConfig(ctx, wasm.Memory, wazero.NewModuleConfig().WithName("env")); err != nil {
		log.Fatal(err)
	}

	args = append([]string{name}, args...)

	cwdAbs, err := filepath.Abs(cwd)
	if err != nil {
		log.Fatal(err)
	}

	var nextTID atomic.Uint64

	// Mount the whole host filesystem and set the guest working directory to the
	// real cwd (see setGuestCwd). Relative, "../", and absolute paths then
	// resolve exactly as they do for the native binary, with no path rewriting.
	cfg := wazero.NewModuleConfig().
		WithSysNanosleep().
		WithSysNanotime().
		WithSysWalltime().
		WithStderr(stderr).
		WithStdout(stdout).
		WithStdin(stdin).
		WithRandSource(rand.Reader).
		WithArgs(args...).
		WithFSConfig(wazero.NewFSConfig().
			WithDirMount(hostMountRoot(cwdAbs), "/"))
	for _, env := range os.Environ() {
		k, v, _ := strings.Cut(env, "=")
		cfg = cfg.WithEnv(k, v)
	}

	if _, err := rt.NewHostModuleBuilder("wasi").NewFunctionBuilder().
		WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, _ api.Module, stack []uint64) {
			tid := nextTID.Add(1)
			startArg := stack[0]
			child, err := rt.InstantiateModule(ctx, code, cfg.
				// Don't need to execute start functions again in child, it crashes anyways.
				WithStartFunctions().
				WithName(""))
			if err != nil {
				panic(err)
			}
			go func() {
				defer func() {
					_ = child.Close(ctx)
				}()
				// wasi_thread_start(thread_id: i32, start_arg: i32)
				if _, err := child.ExportedFunction("wasi_thread_start").Call(ctx, tid, startArg); err != nil {
					log.Printf("wasi_thread_start (tid %d): %v", tid, err)
				}
			}()
			stack[0] = tid
		}), []api.ValueType{api.ValueTypeI32}, []api.ValueType{api.ValueTypeI32}).
		Export("thread-spawn").
		Instantiate(ctx); err != nil {
		log.Fatal(err)
	}

	// Instantiate without running _start to set cwd before the program runs.
	mod, err := rt.InstantiateModule(ctx, code, cfg.WithStartFunctions())
	if err != nil {
		log.Fatal(err)
	}
	if err := setGuestCwd(ctx, mod, cwdAbs); err != nil {
		log.Fatal(err)
	}
	_, err = mod.ExportedFunction("_start").Call(ctx)
	return exitCode(err)
}

// hostMountRoot returns the host directory to mount at the guest root. On Windows,
// this is the volume of cwd. Other volumes aren't accessible but this can be added
// in the future if there is demand.
func hostMountRoot(cwdAbs string) string {
	vol := filepath.VolumeName(cwdAbs)
	if vol == "" {
		return "/"
	}
	return vol + string(filepath.Separator)
}

func setGuestCwd(ctx context.Context, mod api.Module, cwdAbs string) error {
	cwdVar := mod.ExportedGlobal("__wasilibc_cwd") // value is &__wasilibc_cwd
	malloc := mod.ExportedFunction("malloc")
	// Drop Windows volume prefix
	guestCwd := filepath.ToSlash(strings.TrimPrefix(cwdAbs, filepath.VolumeName(cwdAbs)))

	buf := append([]byte(guestCwd), 0)
	res, err := malloc.Call(ctx, uint64(len(buf)))
	if err != nil {
		return fmt.Errorf("malloc cwd buffer: %w", err)
	}
	// wasm32 addresses always fit in uint32.
	strAddr := uint32(res[0]) //nolint:gosec
	if strAddr == 0 {
		return errors.New("malloc returned null for cwd buffer")
	}
	mem := mod.Memory()
	if !mem.Write(strAddr, buf) {
		return fmt.Errorf("write cwd string at %d (%d bytes)", strAddr, len(buf))
	}
	// Store the string's address into the pointer variable: *(&cwd) = strAddr.
	cwdVarAddr := uint32(cwdVar.Get()) //nolint:gosec // wasm32 address fits in uint32
	if !mem.WriteUint32Le(cwdVarAddr, strAddr) {
		return fmt.Errorf("write cwd pointer at %d", cwdVarAddr)
	}
	return nil
}

// exitCode maps a module invocation error to a process exit code: a clean exit
// is 0, a guest proc_exit surfaces its code, and anything else is fatal.
func exitCode(err error) int {
	if err == nil {
		return 0
	}
	sErr := &sys.ExitError{}
	if errors.As(err, &sErr) {
		return int(sErr.ExitCode())
	}
	log.Fatal(err)
	return 0
}
