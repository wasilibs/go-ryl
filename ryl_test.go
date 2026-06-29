package yamllint

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wasilibs/go-ryl/internal/runner"
)

func TestRuns(t *testing.T) {
	stdin := bytes.Buffer{}
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	// We lint the whole folder to verify parallelism.
	ret := runner.Run("ryl", []string{"check", "-d", "extends: default", "-f", "parsable", "testdata"}, &stdin, &stdout, &stderr, ".")
	if want := 1; ret != want {
		t.Fatalf("unexpected return code: have %d, want %d", ret, want)
	}

	wantLines := []string{
		`testdata/error.yaml:1:1: [warning] missing document start "---" (document-start)`,
		`testdata/error.yaml:1:10: [error] too many spaces inside braces (braces)`,
		`testdata/error.yaml:1:27: [error] too many spaces inside braces (braces)`,
		`testdata/test.yaml:1:1: [warning] missing document start "---" (document-start)`,
	}
	for _, line := range wantLines {
		if !strings.Contains(stderr.String(), line) {
			t.Fatalf("missing diagnostic:\n%s\nfull stderr:\n%s", line, stderr.String())
		}
	}
}
