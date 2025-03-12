package scan

import (
	"context"
	"os"
	"path/filepath"
	"slices"

	scalibr "github.com/google/osv-scalibr"
	"github.com/google/osv-scalibr/extractor"
	"github.com/google/osv-scalibr/extractor/filesystem/language/java/javalockfile"
	el "github.com/google/osv-scalibr/extractor/filesystem/list"
	scalibrfs "github.com/google/osv-scalibr/fs"
	scalog "github.com/google/osv-scalibr/log"
	log "github.com/sirupsen/logrus"
)

func isTestDep(i *extractor.Inventory) bool {
	if metadata, ok := i.Metadata.(*javalockfile.Metadata); ok {
		if slices.Contains(metadata.DepGroups(), "test") {
			return true
		}
	}
	return false
}

func findSkipDirs(root string, compare []string) ([]string, error) {
	var results []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a directory
		if !info.IsDir() {
			return nil
		}

		// Check if the directory name (not the full path) matches any string in compare
		if slices.Contains(compare, info.Name()) {
			results = append(results, path)
		}

		return nil
	})

	return results, err
}

func GetJavaDeps(path string, skipDirs []string) (map[string]string, error) {
	log.Debugf("GetJavaDeps %s", path)
	scalog.SetLogger(log.StandardLogger())

	// osv-scalibr prints annoying messages at InfoLevel
	origLevel := log.GetLevel()
	if origLevel >= log.InfoLevel {
		log.SetLevel(log.WarnLevel)
		defer log.SetLevel(origLevel)
	}

	// annoyingly, scalibr requires fullpaths for DirsToSkip, so we need to find these ourselves
	fullSkipPaths, err := findSkipDirs(path, skipDirs)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(path)
	sc := &scalibr.ScanConfig{
		ScanRoots:             scalibrfs.RealFSScanRoots(dir),
		FilesystemExtractors:  el.Java,
		PrintDurationAnalysis: false,
		DirsToSkip:            fullSkipPaths,
	}

	gathered := make(map[string]string, 0)
	results := scalibr.New().Scan(context.Background(), sc)
	for _, i := range results.Inventories {
		if isTestDep(i) {
			log.Debugf("skipping test dependency %s@%s", i.Name, i.Version)
			continue
		}
		p := i.Extractor.ToPURL(i)
		name := p.Name
		if len(p.Namespace) > 0 {
			name = p.Namespace + "/" + p.Name
		}
		// prefer empty string for versions, over "0" or "unknown"
		version := p.Version
		if p.Version == "0" || p.Version == "unknown" {
			version = ""
		}

		gathered[name] = version
	}

	return gathered, nil
}
