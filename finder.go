package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kr/fs"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

// OrderedRanks contains paths ordered by their `Distance` if equal then `Target`.
type OrderedRanks []fuzzy.Rank

func (r OrderedRanks) Len() int {
	return len(r)
}

func (r OrderedRanks) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r OrderedRanks) Less(i, j int) bool {
	if r[i].Distance < r[j].Distance {
		return true
	}
	if r[i].Distance > r[j].Distance {
		return false
	}
	return r[i].Target < r[j].Target
}

// PkgFinder finds a Go package.
type PkgFinder struct {
	// gopath points to $GOPATH/src.
	gopath string

	// cache caches the folder contents for faster lookups.
	cache *cache

	// depthLimit sets the maximum depth of the search.
	// Set it to -1 for infinite depth, otherwise the max depth.
	depthLimit int
}

// NewPkgFinder creates a new `PkgFinder` relative to `path`.
func NewPkgFinder(path string, depth int) *PkgFinder {
	return &PkgFinder{
		gopath:     path,
		cache:      newCache(os.ExpandEnv(CacheFile)),
		depthLimit: depth,
	}
}

// checkRelPath scans every component of the relative path until it finds a
// direct match.
func (w *PkgFinder) checkRelPath(find string, components *[]string) bool {
	for x := len(*components) - 1; x >= 0; x-- {
		if find == filepath.Join((*components)[x:]...) {
			return true
		}
	}
	return false
}

func (w *PkgFinder) walker(root string, find string) []string {
	prevCache := make([]string, 0, 10)
	walker := fs.Walk(root)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		stat := walker.Stat()
		isDir := stat.IsDir()
		path := walker.Path()

		// pkg is the package path, so we only use its parent dir if it is a file
		// otherwise we through off the depth calculation
		var pkg string
		if isDir {
			pkg, _ = filepath.Rel(w.gopath, path)
		} else {
			pkg, _ = filepath.Rel(w.gopath, filepath.Dir(path))
		}

		//   `root` |  0  |  1  | 2
		// ........./$repo/$user/$pkg
		components := strings.Split(pkg, string(filepath.Separator))
		depth := len(components)

		depthLimitHit := w.depthLimit > -1 && depth >= w.depthLimit
		if depthLimitHit {
			walker.SkipDir()
		}
		if path == w.gopath {
			continue
		}
		// Skip if path contains .git or vendor
		if isDir && (strings.HasPrefix(stat.Name(), ".") || strings.HasPrefix(stat.Name(), "_") || strings.Contains(path, "vendor")) {
			walker.SkipDir()
			continue
		}
		if w.cache.fullScan {
			mtime := stat.ModTime().Unix()
			w.cache.cacheStorage.Put(pkg, mtime)
		} else {
			mtime := stat.ModTime().Unix()
			ret, full := w.cache.contains(pkg)
			if !full {
				w.cache.changed = true
				w.cache.cacheStorage.Put(ret, mtime)
			}
			if found := w.checkRelPath(find, &components); found {
				prevCache = append(prevCache, pkg)
				if len(prevCache) > 10 {
					return prevCache
				}
				continue
			}
			// Due to how inodes work (the current inode's mtime only changes if a
			// direct child is modified, it remains unchanged for grandchild and so
			// on) we can only use mtime to skip the last level
			if w.depthLimit-1 == depth {
				v, found := w.cache.cacheStorage.Get(pkg)
				if !found {
					continue
				}
				prevMtime := v.(int64)
				if mtime <= prevMtime {
					walker.SkipDir()
				}
			}
		}
	}
	return prevCache
}

// Find a package by the given key
func (w *PkgFinder) Find(find string) (OrderedRanks, error) {
	// If absolute path then go straight to it
	if path.IsAbs(find) {
		return OrderedRanks{
			{
				Target: find,
			},
		}, nil
	}
	// If path is a real path relative to gopath then use it
	abs := filepath.Join(w.gopath, find)
	_, err := os.Stat(abs)
	if err == nil {
		return OrderedRanks{
			{
				Target: abs,
			},
		}, nil
	}
	defer func() {
		if err := w.cache.save(); err != nil {
			fmt.Fprintln(os.Stderr, "error during cache saving:", err)
		}
	}()
	if w.cache.fullScan {
		_ = w.walker(w.gopath, "")
		w.cache.changed = true
		w.cache.fullScan = false
	}
	// If subpath exists in cache, try it
	paths, fullMatch := w.cache.contains(find)
	if !fullMatch && len(paths) > 0 {
		return w.recheckPaths(paths)
	}
	pkgs := w.walker(w.gopath, find)
	pkgsLen := len(pkgs)
	if pkgs != nil && pkgsLen > 0 {
		if pkgsLen > 1 {
			found := make(OrderedRanks, len(pkgs))
			for i, r := range pkgs {
				found[i] = fuzzy.Rank{
					Target: filepath.Join(w.gopath, r),
				}
			}
			return found, nil
		}
		return OrderedRanks{
			{
				Target: filepath.Join(w.gopath, pkgs[0]),
			},
		}, nil
	}
	// Search in the whole of the hot cache for the best 10 match
	found := w.findSomeMatch(find, 10)

	return found, nil
}

func (w *PkgFinder) recheckPaths(paths []string) (OrderedRanks, error) {
	ret := make(OrderedRanks, 0, len(paths))
	for _, path := range paths {
		var err error
		abs := filepath.Join(w.gopath, path)
		// If we just did a fullScan, cache is valid so no need to check
		if !w.cache.fullScan {
			_, err = os.Stat(abs)
		}
		if err == nil {
			ret = append(ret, fuzzy.Rank{
				Target: abs,
			})
		}
	}
	return ret, nil
}

func (w *PkgFinder) findSomeMatch(find string, num int) OrderedRanks {
	return w.sort(find, w.cache.cacheStorage.Keys(), num)
}

func (w *PkgFinder) sort(by string, paths []interface{}, max int) OrderedRanks {
	matches := make(map[string]fuzzy.Rank)

	for i, p := range paths {
		entry := p.(string)
		path := append(strings.Split(entry, string(filepath.Separator)), entry)
		ranks := fuzzy.RankFindFold(by, path)

		for _, r := range ranks {
			if r.Distance > 10 {
				continue
			}
			m, ok := matches[entry]

			if (ok && r.Distance < m.Distance) || !ok {
				r.Target = filepath.Join(w.gopath, entry)
				matches[entry] = r
				i++
			}
		}
	}
	ret := make(OrderedRanks, len(matches))
	if len(matches) == 0 {
		return ret
	}
	var i int
	for _, r := range matches {
		ret[i] = r
		i++
	}
	sort.Sort(ret)
	if len(ret) > 10 {
		ret = ret[:10]
	}
	return ret
}
