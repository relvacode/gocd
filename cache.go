package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

const (
	// CacheFile is the file used for caching the folder structure.
	CacheFile = "$HOME/.cache/gocd/cache"

	// PrevsFile is the file used for caching the at most 10 folders when that
	// 'fuzzy' matched the queried package.
	PrevsFile = "$HOME/.cache/gocd/prevs"
)

// cache caches the contents of a directory for faster lookups.
type cache struct {
	file     string
	maxDepth int
	cacheStorage

	// changed is set when the folder structure has changed compared to the cache
	changed  bool
	fullScan bool
}

// Creates and loads (according to `file`) a new `cache`.
func newCache(file string) *cache {
	c := &cache{
		file: file,
		cacheStorage: cacheStorage{
			rbt.NewWithStringComparator(),
		},
	}
	c.loadCache()
	return c
}

// load loads the cache file. If it does not exists it creates it.
func (c *cache) loadCache() {
	if _, err := os.Stat(c.file); err != nil {
		c.changed = true
		c.fullScan = true
		return
	}
	if err := c.cacheStorage.deserialize(c.file); err != nil {
		c.loadCacheFail(err)
		return
	}
}

func (c *cache) loadCacheFail(err error) {
	fmt.Fprintln(os.Stderr, err)
	c.changed = true
	c.fullScan = true
}

// save saves the directory structure into the cache file.
func (c *cache) save() error {
	if !c.changed && !c.fullScan {
		return nil
	}
	err := os.MkdirAll(filepath.Dir(c.file), os.ModePerm)
	if err != nil {
		return err
	}
	if err := c.cacheStorage.serialize(c.file); err != nil {
		return err
	}
	return nil
}

func (c *cache) contains(path string) ([]string, bool) {
	ret := make([]string, 0, 10)
	_, found := c.cacheStorage.Get(path)
	if found {
		return ret, true
	}
	it := c.cacheStorage.Iterator()
	for it.Next() {
		entry := it.Key().(string)
		components := strings.Split(entry, string(filepath.Separator))
		for x := len(components) - 1; x >= 0; x-- {
			p := filepath.Join(components[x:]...)
			if path == p {
				ret = append(ret, entry)
				if len(ret) > 10 {
					return ret, false
				}
			}
		}
	}
	return ret, false
}

type cacheStorageEntry struct {
	Path  string
	Mtime int64
}

type cacheStorage struct {
	*rbt.Tree // if no need for ordered keys overkill
}

func (tree *cacheStorage) deserialize(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	elements := make([]cacheStorageEntry, tree.Size())
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&elements); err != nil {
		return err
	}
	tree.Clear()
	for _, elem := range elements {
		tree.Put(elem.Path, elem.Mtime)
	}
	return nil
}

func (tree *cacheStorage) serialize(path string) error {
	elements := make([]cacheStorageEntry, tree.Size())
	it := tree.Iterator()
	for it.Next() {
		elem := cacheStorageEntry{
			Path:  it.Key().(string),
			Mtime: it.Value().(int64),
		}
		elements = append(elements, elem)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err = enc.Encode(elements); err != nil {
		return err
	}
	return nil
}

func loadPrevs() (OrderedRanks, error) {
	prevsFile := os.ExpandEnv(PrevsFile)
	file, err := os.OpenFile(prevsFile, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	prevs := make(OrderedRanks, 10)
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&prevs); err != nil {
		return nil, err
	}
	return prevs, nil
}

func savePrevs(prevs OrderedRanks) error {
	file, err := os.OpenFile(os.ExpandEnv(PrevsFile), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err = enc.Encode(prevs); err != nil {
		return err
	}
	return nil
}
