package example_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	acmev1 "_/internal/pb/acme/v1"

	"github.com/google/go-cmp/cmp"
)

func TestProtoJSON(t *testing.T) {
	t.Run("acme.v1.Example", func(t *testing.T) {
		t.Parallel()
		EachFile[acmev1.Example](t, "testdata/acme.v1.Example")
	})

	t.Run("acme.v1.SomethingElse", func(t *testing.T) {
		t.Parallel()
		EachFile[acmev1.SomethingElse](t, "testdata/acme.v1.SomethingElse")
	})
}

func EachFile[P any](t *testing.T, baseDir string) {
	t.Helper()
	files, err := os.ReadDir(baseDir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %v", baseDir, err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()
			src, err := os.ReadFile(filepath.Join(baseDir, file.Name()))
			if err != nil {
				t.Fatalf("failed to read file %s: %v", file.Name(), err)
			}
			var value P
			if err := json.Unmarshal(src, &value); err != nil {
				t.Fatalf("failed to unmarshal JSON: %v", err)
			}
			t.Logf("Parsed value: %+v", value)

			dst, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				t.Fatalf("failed to marshal JSON: %v", err)
			}

			// Editors typically have a trailing newline, so add it to the output for the sake of test.
			dst = append(dst, '\n')

			if diff := cmp.Diff(src, dst); diff != "" {
				t.Errorf("JSON does not match original %q (-want +got):\n%s", file.Name(), diff)
			}
		})
	}
}
