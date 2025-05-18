package unzip_test

import (
	"archive/zip"
	"os"
	"path/filepath"
	"pb_launcher/helpers/unzip"
	"strings"
	"testing"
)

func createTestZip(t *testing.T, files map[string]string) (string, func()) {
	t.Helper()

	// Crear archivo temporal
	tempFile, err := os.CreateTemp("", "test-*.zip")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	zipWriter := zip.NewWriter(tempFile)
	for name, content := range files {
		writer, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("failed to create zip entry %s: %v", name, err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write to zip entry %s: %v", name, err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatalf("failed to close zip writer: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Devolver la ruta y la función cleanup
	return tempFile.Name(), func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("failed to clean up temp file %s: %v", tempFile.Name(), err)
		}
	}
}

func TestExtract(t *testing.T) {
	testCases := []struct {
		name          string
		files         map[string]string
		expectedFiles []string
		expectError   bool
	}{
		{
			name:          "empty zip file",
			files:         map[string]string{},
			expectedFiles: []string{},
			expectError:   false,
		},
		{
			name: "single large file",
			files: map[string]string{
				"large_file.txt": strings.Repeat("A", 10*1024*1024), // 10 MB
			},
			expectedFiles: []string{"large_file.txt"},
			expectError:   false,
		},
		{
			name: "files with special characters",
			files: map[string]string{
				"spécial-chär.txt": "some content",
				"sp@ce fil&e.txt":  "with special chars",
			},
			expectedFiles: []string{"spécial-chär.txt", "sp@ce fil&e.txt"},
			expectError:   false,
		},
		{
			name: "deeply nested directories",
			files: map[string]string{
				"nested/dir/with/many/levels/file.txt": "deep content",
			},
			expectedFiles: []string{"nested/dir/with/many/levels/file.txt"},
			expectError:   false,
		},
		{
			name: "hidden files",
			files: map[string]string{
				".hidden_file":             "hidden content",
				".hidden/dir/.hidden_file": "nested hidden file",
			},
			expectedFiles: []string{".hidden_file", ".hidden/dir/.hidden_file"},
			expectError:   false,
		},
		{
			name: "symlink traversal attempt",
			files: map[string]string{
				"../../../../etc/passwd": "malicious content",
			},
			expectError: true,
		},
		{
			name: "directory with trailing slash",
			files: map[string]string{
				"dir_with_slash/":         "",
				"dir_with_slash/file.txt": "content",
			},
			expectedFiles: []string{"dir_with_slash/", "dir_with_slash/file.txt"},
			expectError:   false,
		},
		{
			name: "files with no extension",
			files: map[string]string{
				"no_extension": "no extension content",
			},
			expectedFiles: []string{"no_extension"},
			expectError:   false,
		},
		{
			name: "empty directory",
			files: map[string]string{
				"empty_dir/": "",
			},
			expectedFiles: []string{"empty_dir/"},
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zipPath, cleanup := createTestZip(t, tc.files)
			if cleanup != nil {
				defer cleanup()
			}

			dest, err := os.MkdirTemp("", "unzip-test-*")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(dest)

			u := unzip.NewUnzip()
			files, err := u.Extract(zipPath, dest)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			for _, expectedFile := range tc.expectedFiles {
				fullPath := filepath.Join(dest, expectedFile)
				if _, err := os.Stat(fullPath); err != nil {
					t.Errorf("expected file %s not found: %v", fullPath, err)
				}
			}

			if len(files) != len(tc.expectedFiles) {
				t.Errorf("expected %d files but got %d", len(tc.expectedFiles), len(files))
			}
		})
	}
}
