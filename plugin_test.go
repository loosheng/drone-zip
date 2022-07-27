package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
