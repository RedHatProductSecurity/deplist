package deplist

import (
	"fmt"
)

// Bitmask type allows easy tagging of what langs there are
type Bitmask uint32

// Dependency per dependency info
type Dependency struct {
	DepType Bitmask  // golang, nodejs, python etc
	Path    string   // the module path, github.com/teris-io/shortid
	Version string   // v0.0.0-20171029131806-771a37caa5cf
	Files   []string // if available, list of all files for a package
	// /usr/lib/go-1.13/src/regexp/syntax/compile.go
	// /usr/lib/go-1.13/src/regexp/syntax/doc.go
}

// DepFoundAddFlag add a lang type to the bitmask
func (f *Bitmask) DepFoundAddFlag(flag Bitmask) { *f |= flag }

// DepFoundHasFlag deteremine if bitmask has a lang type
func (f Bitmask) DepFoundHasFlag(flag Bitmask) bool { return f&flag != 0 }

func (d *Dependency) ToString() string {
	return fmt.Sprintf("pkg:%s/%s@%s", GetLanguageStr(d.DepType), d.Path, d.Version)
}

// GetLanguageStr returns from a bitmask return the ecosystem name
func GetLanguageStr(bm Bitmask) string {
	if bm&LangGolang != 0 {
		return "go"
	} else if bm&LangJava != 0 {
		return "mvn"
	} else if bm&LangNodeJS != 0 {
		return "npm"
	} else if bm&LangPython != 0 {
		return "pypi"
	} else if bm&LangRuby != 0 {
		return "gem"
	}
	return "unknown"
}
