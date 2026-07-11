package fsutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	// ErrPathTraversal is returned when path traversal or root escape is detected.
	ErrPathTraversal = errors.New("path traversal detected")
	// ErrSymlink is returned when a symlink is detected.
	ErrSymlink = errors.New("symlinks are not allowed")
	// ErrHardlink is returned when a hardlink is detected.
	ErrHardlink = errors.New("hardlinks are not allowed")
	// ErrDuplicate is returned when a duplicate name is detected.
	ErrDuplicate = errors.New("item already exists")
	// ErrInvalidName is returned when a folder or file name is invalid.
	ErrInvalidName = errors.New("invalid item name")
)

// ValidatePath checks if the requested path is safe and inside the root directory.
// It returns the cleaned absolute path to the target.
func ValidatePath(root, requestedPath string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path of root: %w", err)
	}

	// Reject absolute paths or those beginning with a path separator
	if filepath.IsAbs(requestedPath) ||
		strings.HasPrefix(requestedPath, "/") ||
		strings.HasPrefix(requestedPath, "\\") {
		return "", ErrPathTraversal
	}

	// Join root and requestedPath, then clean it.
	targetPath := filepath.Clean(filepath.Join(absRoot, requestedPath))

	// Get relative path between root and target to check for traversal.
	rel, err := filepath.Rel(absRoot, targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to compute relative path: %w", err)
	}

	// If relative path starts with ".." or is absolute escaping root, it's a traversal.
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ErrPathTraversal
	}

	return targetPath, nil
}

// RejectSymlinksAndHardlinks checks if the target path is a symlink or hardlink.
// For existing files/directories, it returns an error if a symlink or hardlink is detected.
func RejectSymlinksAndHardlinks(path string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// Check if symlink
	if fi.Mode()&os.ModeSymlink != 0 {
		return ErrSymlink
	}

	// Check if hardlink (only for regular files, directories can have Nlink > 1)
	if !fi.IsDir() {
		if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
			if stat.Nlink > 1 {
				return ErrHardlink
			}
		}
	}

	return nil
}

// CheckDuplicate checks if a file or directory with the given name already exists inside dirPath.
func CheckDuplicate(dirPath, name string) error {
	target := filepath.Join(dirPath, name)
	if _, err := os.Lstat(target); err == nil {
		return ErrDuplicate
	}
	return nil
}

// ValidateName ensures that the name does not contain path separators, null bytes,
// and is not empty, "." or "..".
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidName)
	}
	if strings.ContainsAny(name, "/\\") || strings.Contains(name, "\x00") {
		return fmt.Errorf("%w: name cannot contain path separators or null bytes", ErrInvalidName)
	}
	if name == "." || name == ".." {
		return fmt.Errorf("%w: name cannot be '.' or '..'", ErrInvalidName)
	}
	if len(name) > 255 {
		return fmt.Errorf("%w: name too long", ErrInvalidName)
	}
	return nil
}
