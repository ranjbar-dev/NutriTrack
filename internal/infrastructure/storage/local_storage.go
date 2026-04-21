package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// LocalStorage saves files to the local filesystem.
type LocalStorage struct {
	BasePath string // e.g. "uploads"
	BaseURL  string // e.g. "/uploads" or "https://..."
}

// NewLocalStorage creates a new LocalStorage with the given base path and URL prefix.
func NewLocalStorage(basePath, baseURL string) *LocalStorage {
	return &LocalStorage{BasePath: basePath, BaseURL: baseURL}
}

// SaveAvatar saves a file to <basePath>/avatars/<uuid>.<ext> and returns the URL path.
func (s *LocalStorage) SaveAvatar(src io.Reader, ext string) (string, error) {
	dir := filepath.Join(s.BasePath, "avatars")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.%s", uuid.NewString(), ext)
	fullPath := filepath.Join(dir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/avatars/%s", s.BaseURL, filename), nil
}
