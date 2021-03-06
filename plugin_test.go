package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin(t *testing.T) {

	t.Run("zip file ./test/*", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Logf("tmpDir: %v", tmpDir)
		p := Plugin{
			Input:  []string{"./test/*"},
			Output: filepath.Join(tmpDir, "test.zip"),
		}
		p.Exec()

		assert.Equal(
			t,
			[]string{
				"test/a.txt",
				"test/b.js",
			},
			zipFiles(t, filepath.Join(tmpDir, "test.zip")),
		)

	})

	t.Run("zip file ./test", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Logf("tmpDir: %v", tmpDir)
		p := Plugin{
			Input:  []string{"./test"},
			Output: filepath.Join(tmpDir, "dot-test.zip"),
		}
		p.Exec()

		assert.Equal(
			t,
			[]string{
				"test/a.txt",
				"test/b.js",
				"test/foo/a.js",
				"test/foo/b.txt",
				"test/foo/bar/bar.js"},
			zipFiles(t, filepath.Join(tmpDir, "dot-test.zip")),
		)

	})

	t.Run("zip file ./test/**/*", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Logf("tmpDir: %v", tmpDir)
		p := Plugin{
			Input:  []string{"./test/**/*"},
			Output: filepath.Join(tmpDir, "db-star-test.zip"),
		}
		p.Exec()

		assert.Equal(
			t,
			[]string{
				"test/a.txt",
				"test/b.js",
				"test/foo/a.js",
				"test/foo/b.txt",
				"test/foo/bar/bar.js",
			},
			zipFiles(t, filepath.Join(tmpDir, "db-star-test.zip")),
		)

	})

	t.Run("zip file ./test/**/*.js", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Logf("tmpDir: %v", tmpDir)
		p := Plugin{
			Input:  []string{"./test/**/*.js"},
			Output: filepath.Join(tmpDir, "db-star-js-test.zip"),
		}
		p.Exec()

		assert.Equal(
			t,
			[]string{
				"test/b.js",
				"test/foo/a.js",
				"test/foo/bar/bar.js",
			},
			zipFiles(t, filepath.Join(tmpDir, "db-star-js-test.zip")),
		)
	})

	t.Run("zip file ./test/**/*.txt", func(t *testing.T) {
		tmpDir := t.TempDir()
		t.Logf("tmpDir: %v", tmpDir)
		p := Plugin{
			Input:  []string{"./test/**/*.txt"},
			Output: filepath.Join(tmpDir, "db-star-txt-test.zip"),
		}
		p.Exec()

		assert.Equal(
			t,
			[]string{
				"test/a.txt",
				"test/foo/b.txt",
			},
			zipFiles(t, filepath.Join(tmpDir, "db-star-txt-test.zip")),
		)
	})

}

func zipFiles(t *testing.T, path string) []string {
	t.Helper()
	f, err := os.Open(path)
	assert.NoError(t, err)
	info, err := f.Stat()
	assert.NoError(t, err)
	r, err := zip.NewReader(f, info.Size())
	assert.NoError(t, err)
	paths := make([]string, len(r.File))
	for i, zf := range r.File {
		paths[i] = zf.Name
	}
	return paths
}

func TestIsDir(t *testing.T) {

	t.Run("It's a folder", func(t *testing.T) {
		assert.True(t, IsDir("./test/"))
		assert.True(t, IsDir("test/"))
		assert.True(t, IsDir("test"))

	})

	t.Run("It's not a folder", func(t *testing.T) {
		assert.False(t, IsDir("./main.go"))
		assert.False(t, IsDir("main.go"))
		assert.False(t, IsDir("aaa"))
	})

}

func TestContains(t *testing.T) {

	t.Run("string slice contains", func(t *testing.T) {
		assert.True(t, Contains([]string{"a", "b"}, "a"))
		assert.True(t, Contains([]string{"a", "b"}, "b"))

		assert.True(t, Contains([]string{"ab"}, "b"))
		assert.True(t, Contains([]string{"ab"}, "a"))

		assert.True(t, Contains([]string{"This is a sentence"}, "sentence"))
		assert.True(t, Contains([]string{"This is a sentence"}, "sent"))
	})

	t.Run("string slice does not contains", func(t *testing.T) {
		assert.False(t, Contains([]string{"a", "b"}, "c"))
		assert.False(t, Contains([]string{"a", "b"}, "d"))
	})
}

func TestGetFilePaths(t *testing.T) {
	t.Run("glob match", func(t *testing.T) {

		assert.Equal(t, []string{"test/a.txt"}, getFilePaths("./test/a.txt"))

		assert.Equal(t, []string{"test/a.txt", "test/b.js"}, getFilePaths("./test/*"))

		assert.Equal(t,
			[]string{
				"test/a.txt",
				"test/b.js",
				"test/foo/a.js",
				"test/foo/b.txt",
				"test/foo/bar/bar.js"},
			getFilePaths("./test"),
		)

		assert.Equal(t,
			[]string{
				"test/a.txt",
				"test/b.js",
				"test/foo/a.js",
				"test/foo/b.txt",
				"test/foo/bar/bar.js"},
			getFilePaths("./test/**"),
		)

		assert.Equal(t,
			[]string{
				"test/a.txt",
				"test/b.js",
				"test/foo/a.js",
				"test/foo/b.txt",
				"test/foo/bar/bar.js",
			},
			getFilePaths("./test/**/*"),
		)

		assert.Equal(t,
			[]string{
				"test/a.txt",
			},
			getFilePaths("./test/*.txt"),
		)

		assert.Equal(t,
			[]string{
				"test/a.txt",
				"test/foo/b.txt",
			},
			getFilePaths("./test/**/*.txt"),
		)

		assert.Equal(t,
			[]string{
				"test/b.js",
			},
			getFilePaths("./test/*.js"),
		)

		assert.Equal(t,
			[]string{
				"test/b.js",
				"test/foo/a.js",
				"test/foo/bar/bar.js",
			},
			getFilePaths("./test/**/*.js"),
		)

	})
}
