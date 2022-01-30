package mod

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/modcache"
)

// ErrCantInfer occurs when you can't infer the module path from the $GOPATH.
var ErrCantInfer = errors.New("mod: unable to infer the module path")

// ErrFileNotFound occurs when no go.mod can be found
var ErrFileNotFound = fmt.Errorf("unable to find go.mod: %w", fs.ErrNotExist)

type Option = func(f *Finder)

// New modfile loader
func New(options ...Option) *Finder {
	finder := &Finder{
		cache: modcache.Default(),
		fsys:  osfs{},
	}
	for _, option := range options {
		option(finder)
	}
	return finder
}

// WithCache uses a custom mod cache instead of the default
func WithCache(cache *modcache.Cache) func(f *Finder) {
	return func(f *Finder) {
		f.cache = cache
	}
}

// WithFS specifies an filesystem to use while finding
func WithFS(fsys fs.FS) func(f *Finder) {
	return func(f *Finder) {
		f.fsys = fsys
	}
}

// FindDirectory traverses up the filesystem until it finds a directory
// containing go.mod or returns an error trying.
func FindDirectory(dir string) (abs string, err error) {
	dir, err = findDirectory(osfs{}, dir)
	if err != nil {
		return "", err
	}
	return filepath.Abs(dir)
}

func findDirectory(fsys fs.FS, dir string) (abs string, err error) {
	path := filepath.Join(dir, "go.mod")
	// Check if this path exists, otherwise recursively traverse towards root
	if _, err = fs.Stat(fsys, path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		nextDir := filepath.Dir(dir)
		if nextDir == dir {
			return "", ErrFileNotFound
		}
		return findDirectory(fsys, filepath.Dir(dir))
	}
	return dir, nil
}