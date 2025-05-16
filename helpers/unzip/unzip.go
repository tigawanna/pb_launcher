package unzip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Unzip struct{}

func NewUnzip() *Unzip {
	return &Unzip{}
}

func (uz Unzip) Extract(source, destination string) ([]string, error) {
	r, err := zip.OpenReader(source)
	if err != nil {
		return nil, fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(destination, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %w", err)
	}

	var extractedFiles []string
	for _, f := range r.File {
		if err := uz.extractAndWriteFile(destination, f); err != nil {
			return nil, fmt.Errorf("failed to extract file %s: %w", f.Name, err)
		}
		extractedFiles = append(extractedFiles, f.Name)
	}

	return extractedFiles, nil
}

func (Unzip) extractAndWriteFile(destination string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s in zip archive: %w", f.Name, err)
	}
	defer rc.Close()

	cleanDest, err := filepath.Abs(destination)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute destination path: %w", err)
	}

	targetPath := filepath.Join(cleanDest, f.Name)
	relPath, err := filepath.Rel(cleanDest, targetPath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("illegal file path detected: %s", targetPath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories: %w", err)
		}

		dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, rc); err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
	}
	return nil
}
