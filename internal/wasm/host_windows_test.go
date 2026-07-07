//go:build windows

package wasm2go

import "testing"

func TestHostPathToGuestPathWindowsDrive(t *testing.T) {
	got := hostPathToGuestPath(`C:\Users\runneradmin\AppData\Local\Temp\ryl`)
	want := `/C:/Users/runneradmin/AppData/Local/Temp/ryl`
	if got != want {
		t.Fatalf("hostPathToGuestPath() = %q, want %q", got, want)
	}
}

func TestGuestPathToHostPathWindowsDrive(t *testing.T) {
	got := guestPathToHostPath(`\`, `/C:/Users/runneradmin/AppData/Local/Temp/ryl/shell.md`)
	want := `C:\Users\runneradmin\AppData\Local\Temp\ryl\shell.md`
	if got != want {
		t.Fatalf("guestPathToHostPath() = %q, want %q", got, want)
	}
}
