package fsutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestValidatePath(t *testing.T) {
	tmp := t.TempDir()

	// Create some dirs
	sub := filepath.Join(tmp, "sub")
	if err := os.Mkdir(sub, 0700); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		reqPath string
		wantErr error
	}{
		{"valid root", "", nil},
		{"valid sub", "sub", nil},
		{"valid sub with slash", "sub/file.txt", nil},
		{"traversal basic", "..", ErrPathTraversal},
		{"traversal relative dotdot", "sub/../../other", ErrPathTraversal},
		{"traversal absolute style", "/etc/passwd", ErrPathTraversal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidatePath(tmp, tt.reqPath)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePath() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestRejectSymlinksAndHardlinks(t *testing.T) {
	tmp := t.TempDir()

	file := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}

	// Test regular file is okay
	if err := RejectSymlinksAndHardlinks(file); err != nil {
		t.Errorf("expected regular file to be accepted, got: %v", err)
	}

	// Symlinks test (Unix only as Windows needs admin/developer mode for symlinks)
	if runtime.GOOS != "windows" {
		sym := filepath.Join(tmp, "sym.txt")
		if err := os.Symlink(file, sym); err != nil {
			t.Fatal(err)
		}

		if err := RejectSymlinksAndHardlinks(sym); err != ErrSymlink {
			t.Errorf("expected ErrSymlink, got: %v", err)
		}

		// Hardlinks test
		hl := filepath.Join(tmp, "hl.txt")
		if err := os.Link(file, hl); err != nil {
			t.Fatal(err)
		}

		if err := RejectSymlinksAndHardlinks(hl); err != ErrHardlink {
			t.Errorf("expected ErrHardlink on hl, got: %v", err)
		}

		if err := RejectSymlinksAndHardlinks(file); err != ErrHardlink {
			t.Errorf("expected ErrHardlink on original file (nlink > 1), got: %v", err)
		}
	}
}

func TestCheckDuplicate(t *testing.T) {
	tmp := t.TempDir()

	file := filepath.Join(tmp, "file.txt")
	if err := os.WriteFile(file, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := CheckDuplicate(tmp, "file.txt"); err != ErrDuplicate {
		t.Errorf("expected ErrDuplicate, got: %v", err)
	}

	if err := CheckDuplicate(tmp, "other.txt"); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "hello.txt", false},
		{"empty", "", true},
		{"spaces only", "   ", true},
		{"slash", "hello/world", true},
		{"backslash", "hello\\world", true},
		{"dot", ".", true},
		{"dotdot", "..", true},
		{"too long", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearchAndListDir(t *testing.T) {
	tmp := t.TempDir()

	// Set up hierarchy
	// tmp/
	//   file1.txt
	//   dir1/
	//     file2.txt
	//     file_match.log
	//   dir2/
	//     match_file.txt

	if err := os.WriteFile(filepath.Join(tmp, "file1.txt"), []byte("file1"), 0600); err != nil {
		t.Fatal(err)
	}

	dir1 := filepath.Join(tmp, "dir1")
	if err := os.Mkdir(dir1, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "file2.txt"), []byte("file2"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "file_match.log"), []byte("match"), 0600); err != nil {
		t.Fatal(err)
	}

	dir2 := filepath.Join(tmp, "dir2")
	if err := os.Mkdir(dir2, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "match_file.txt"), []byte("match"), 0600); err != nil {
		t.Fatal(err)
	}

	// Test ListDir
	list, err := ListDir(tmp, tmp)
	if err != nil {
		t.Fatalf("ListDir failed: %v", err)
	}
	if len(list) != 3 {
		t.Errorf("expected 3 items in ListDir, got %d", len(list))
	}

	// Test Search
	searchRes, err := Search(tmp, tmp, "match")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(searchRes) != 2 {
		t.Errorf("expected 2 search results, got %d", len(searchRes))
	}

	// Verify paths are relative
	for _, item := range searchRes {
		if filepath.IsAbs(item.Path) {
			t.Errorf("expected relative path, got: %s", item.Path)
		}
		if !strings.Contains(strings.ToLower(item.Name), "match") {
			t.Errorf("expected name to contain 'match', got: %s", item.Name)
		}
	}
}
