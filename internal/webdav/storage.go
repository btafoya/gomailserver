package webdav

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ResourceInfo contains metadata about a WebDAV resource
type ResourceInfo struct {
	Path         string
	Name         string
	IsCollection bool
	ContentType  string
	ContentLen   int64
	ETag         string
	ModTime      time.Time
	CreateTime   time.Time
	ResourceKind string // "collection", "calendar", "addressbook", "principal", "event", "contact", "file"
}

// Storage defines the interface for WebDAV resource storage
type Storage interface {
	// GetResourceInfo retrieves metadata for a resource
	GetResourceInfo(path string) (*ResourceInfo, error)

	// ListChildren returns children of a collection
	ListChildren(path string) ([]*ResourceInfo, error)

	// CreateCollection creates a new collection (directory)
	CreateCollection(path string) error

	// DeleteResource deletes a resource or collection
	DeleteResource(path string) error

	// CopyResource copies a resource to a new location
	CopyResource(src, dst string, overwrite bool) error

	// MoveResource moves a resource to a new location
	MoveResource(src, dst string, overwrite bool) error

	// ReadResource reads the content of a resource
	ReadResource(path string) (io.ReadCloser, error)

	// WriteResource writes content to a resource
	WriteResource(path string, content io.Reader) error

	// Exists checks if a resource exists
	Exists(path string) bool
}

// FileSystemStorage implements Storage using the local filesystem
type FileSystemStorage struct {
	basePath string
	logger   *zap.Logger
}

// NewFileSystemStorage creates a new filesystem-based storage
func NewFileSystemStorage(basePath string, logger *zap.Logger) *FileSystemStorage {
	return &FileSystemStorage{
		basePath: basePath,
		logger:   logger,
	}
}

// resolvePath converts a WebDAV path to a filesystem path
func (s *FileSystemStorage) resolvePath(path string) string {
	// Clean and normalize the path
	cleaned := filepath.Clean(path)
	if cleaned == "." || cleaned == "/" {
		return s.basePath
	}
	// Remove leading slash
	cleaned = strings.TrimPrefix(cleaned, "/")
	return filepath.Join(s.basePath, cleaned)
}

// GetResourceInfo retrieves metadata for a resource
func (s *FileSystemStorage) GetResourceInfo(path string) (*ResourceInfo, error) {
	fsPath := s.resolvePath(path)

	info, err := os.Stat(fsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("resource not found: %s", path)
		}
		return nil, fmt.Errorf("failed to stat resource: %w", err)
	}

	resourceInfo := &ResourceInfo{
		Path:         path,
		Name:         info.Name(),
		IsCollection: info.IsDir(),
		ModTime:      info.ModTime(),
		CreateTime:   info.ModTime(), // Use mod time as fallback
	}

	if !info.IsDir() {
		resourceInfo.ContentLen = info.Size()
		resourceInfo.ContentType = s.detectContentType(path)
		resourceInfo.ETag = s.generateETag(fsPath, info)
	}

	resourceInfo.ResourceKind = s.detectResourceKind(path, info.IsDir())

	return resourceInfo, nil
}

// ListChildren returns children of a collection
func (s *FileSystemStorage) ListChildren(path string) ([]*ResourceInfo, error) {
	fsPath := s.resolvePath(path)

	entries, err := os.ReadDir(fsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("collection not found: %s", path)
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var children []*ResourceInfo
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childInfo, err := s.GetResourceInfo(childPath)
		if err != nil {
			s.logger.Warn("failed to get child info", zap.String("path", childPath), zap.Error(err))
			continue
		}
		children = append(children, childInfo)
	}

	return children, nil
}

// CreateCollection creates a new collection (directory)
func (s *FileSystemStorage) CreateCollection(path string) error {
	fsPath := s.resolvePath(path)

	// Check if parent exists
	parent := filepath.Dir(fsPath)
	if _, err := os.Stat(parent); os.IsNotExist(err) {
		return fmt.Errorf("parent collection does not exist")
	}

	// Check if resource already exists
	if _, err := os.Stat(fsPath); err == nil {
		return fmt.Errorf("resource already exists: %s", path)
	}

	if err := os.Mkdir(fsPath, 0755); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	s.logger.Info("created collection", zap.String("path", path))
	return nil
}

// DeleteResource deletes a resource or collection
func (s *FileSystemStorage) DeleteResource(path string) error {
	fsPath := s.resolvePath(path)

	info, err := os.Stat(fsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("resource not found: %s", path)
		}
		return fmt.Errorf("failed to stat resource: %w", err)
	}

	if info.IsDir() {
		if err := os.RemoveAll(fsPath); err != nil {
			return fmt.Errorf("failed to delete collection: %w", err)
		}
	} else {
		if err := os.Remove(fsPath); err != nil {
			return fmt.Errorf("failed to delete resource: %w", err)
		}
	}

	s.logger.Info("deleted resource", zap.String("path", path))
	return nil
}

// CopyResource copies a resource to a new location
func (s *FileSystemStorage) CopyResource(src, dst string, overwrite bool) error {
	srcPath := s.resolvePath(src)
	dstPath := s.resolvePath(dst)

	// Check if source exists
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("source not found: %s", src)
	}

	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil && !overwrite {
		return fmt.Errorf("destination already exists: %s", dst)
	}

	if srcInfo.IsDir() {
		return s.copyDir(srcPath, dstPath)
	}
	return s.copyFile(srcPath, dstPath)
}

// MoveResource moves a resource to a new location
func (s *FileSystemStorage) MoveResource(src, dst string, overwrite bool) error {
	srcPath := s.resolvePath(src)
	dstPath := s.resolvePath(dst)

	// Check if destination exists
	if _, err := os.Stat(dstPath); err == nil && !overwrite {
		return fmt.Errorf("destination already exists: %s", dst)
	}

	// Remove destination if it exists and overwrite is true
	if overwrite {
		_ = os.RemoveAll(dstPath)
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		return fmt.Errorf("failed to move resource: %w", err)
	}

	s.logger.Info("moved resource", zap.String("src", src), zap.String("dst", dst))
	return nil
}

// ReadResource reads the content of a resource
func (s *FileSystemStorage) ReadResource(path string) (io.ReadCloser, error) {
	fsPath := s.resolvePath(path)

	file, err := os.Open(fsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("resource not found: %s", path)
		}
		return nil, fmt.Errorf("failed to open resource: %w", err)
	}

	return file, nil
}

// WriteResource writes content to a resource
func (s *FileSystemStorage) WriteResource(path string, content io.Reader) error {
	fsPath := s.resolvePath(path)

	// Ensure parent directory exists
	parent := filepath.Dir(fsPath)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	file, err := os.Create(fsPath)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, content); err != nil {
		return fmt.Errorf("failed to write content: %w", err)
	}

	s.logger.Info("wrote resource", zap.String("path", path))
	return nil
}

// Exists checks if a resource exists
func (s *FileSystemStorage) Exists(path string) bool {
	fsPath := s.resolvePath(path)
	_, err := os.Stat(fsPath)
	return err == nil
}

// detectContentType determines the MIME type based on file extension
func (s *FileSystemStorage) detectContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".ics":
		return "text/calendar; charset=utf-8"
	case ".vcf":
		return "text/vcard; charset=utf-8"
	case ".txt":
		return "text/plain; charset=utf-8"
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".xml":
		return "application/xml; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".pdf":
		return "application/pdf"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

// generateETag generates an ETag for a file based on content hash and mod time
func (s *FileSystemStorage) generateETag(fsPath string, info os.FileInfo) string {
	// For performance, use size + mod time instead of content hash for large files
	if info.Size() > 10*1024*1024 { // 10MB threshold
		data := fmt.Sprintf("%s-%d-%d", fsPath, info.Size(), info.ModTime().UnixNano())
		hash := sha256.Sum256([]byte(data))
		return `"` + hex.EncodeToString(hash[:16]) + `"`
	}

	// For smaller files, hash the content
	file, err := os.Open(fsPath)
	if err != nil {
		// Fallback to size + mod time
		data := fmt.Sprintf("%s-%d-%d", fsPath, info.Size(), info.ModTime().UnixNano())
		hash := sha256.Sum256([]byte(data))
		return `"` + hex.EncodeToString(hash[:16]) + `"`
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		data := fmt.Sprintf("%s-%d-%d", fsPath, info.Size(), info.ModTime().UnixNano())
		h := sha256.Sum256([]byte(data))
		return `"` + hex.EncodeToString(h[:16]) + `"`
	}

	return `"` + hex.EncodeToString(hash.Sum(nil)[:16]) + `"`
}

// detectResourceKind determines the WebDAV resource type based on path
func (s *FileSystemStorage) detectResourceKind(path string, isDir bool) string {
	if strings.HasPrefix(path, "/principals/") {
		return "principal"
	}
	if strings.Contains(path, "/calendars/") {
		if isDir {
			return "calendar"
		}
		if strings.HasSuffix(path, ".ics") {
			return "event"
		}
		return "collection"
	}
	if strings.Contains(path, "/addressbooks/") {
		if isDir {
			return "addressbook"
		}
		if strings.HasSuffix(path, ".vcf") {
			return "contact"
		}
		return "collection"
	}
	if isDir {
		return "collection"
	}
	return "file"
}

// copyFile copies a single file
func (s *FileSystemStorage) copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Preserve mod time
	srcInfo, err := os.Stat(src)
	if err == nil {
		_ = os.Chtimes(dst, srcInfo.ModTime(), srcInfo.ModTime())
	}

	return nil
}

// copyDir recursively copies a directory
func (s *FileSystemStorage) copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := s.copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := s.copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
