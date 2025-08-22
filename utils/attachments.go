package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const TempDirPrefix = "emailsvc-attach-"

// SaveBase64Attachments writes attachments to temp files and returns the full paths.
// Caller must call CleanupFiles on returned paths.
func SaveBase64Attachments(atts []struct {
	Filename string `json:"filename" validate:"required"`
	Content  string `json:"content"  validate:"required"`
}) ([]string, error) {
	paths := make([]string, 0, len(atts))
	for i, a := range atts {
		data, err := base64.StdEncoding.DecodeString(a.Content)
		if err != nil {
			cleanupFiles(paths)
			return nil, fmt.Errorf("attachment %d base64 decode failed: %w", i, err)
		}

		name := filepath.Base(a.Filename)
		if name == "" || name == "." {
			name = "attachment-" + strconv.Itoa(i)
		}

		tmpDir, err := os.MkdirTemp("", TempDirPrefix)
		if err != nil {
			cleanupFiles(paths)
			return nil, fmt.Errorf("create temp dir failed: %w", err)
		}

		tmpPath := filepath.Join(tmpDir, name)
		f, err := os.Create(tmpPath)
		if err != nil {
			_ = os.RemoveAll(tmpDir)
			cleanupFiles(paths)
			return nil, fmt.Errorf("create temp file failed: %w", err)
		}
		if _, err := io.Copy(f, bytesReader(data)); err != nil {
			f.Close()
			_ = os.RemoveAll(tmpDir)
			cleanupFiles(paths)
			return nil, fmt.Errorf("write temp file failed: %w", err)
		}
		_ = f.Close()
		paths = append(paths, tmpPath)
	}
	return paths, nil
}

func cleanupFiles(paths []string) {
	for _, p := range paths {
		_ = os.Remove(p)
		dir := filepath.Dir(p)
		_ = os.RemoveAll(dir)
	}
}

// exported helper for callers
func CleanupFiles(paths []string) {
	cleanupFiles(paths)
}

// small bytes reader
func bytesReader(b []byte) io.Reader {
	return &byteReader{b: b}
}

type byteReader struct{ b []byte }

func (r *byteReader) Read(p []byte) (int, error) {
	if len(r.b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, r.b)
	r.b = r.b[n:]
	return n, nil
}
