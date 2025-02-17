package deplist

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/RedHatProductSecurity/deplist/internal/scan"
	"github.com/RedHatProductSecurity/deplist/internal/utils"

	log "github.com/sirupsen/logrus"
)

const LangGeneric = 0

// enums start at 1 to allow us to specify found languages 0 = nil
const (
	LangGolang = 1 << iota
	LangJava
	LangNodeJS
	LangPython
	LangRuby
	LangRust
)

func init() {
	// check for the library required binaries
	languages := map[string]string{
		"yarn":   "yarn",
		"npm":    "npm",
		"go":     "go",
		"bundle": "bundler gem",
	}

	for lang_bin, lang_name := range languages {
		if _, err := exec.LookPath(lang_bin); err != nil {
			log.Fatal(lang_name, " is required in PATH")
		}
	}
}

type Discovered struct {
	deps       []Dependency
	foundTypes Bitmask
}

func addPackagesToDeps(discovered Discovered, pkgs map[string]string, lang Bitmask) Discovered {
	if len(pkgs) > 0 {
		discovered.foundTypes.DepFoundAddFlag(lang)
	}

	for name, version := range pkgs {
		discovered.deps = append(discovered.deps,
			Dependency{
				DepType: lang,
				Path:    strings.TrimSuffix(name, "\n"),
				Version: strings.Replace(version, "v", "", 1),
				Files:   []string{},
			})
	}
	return discovered
}

var defaultIgnore []string = []string{
	".git",
	"docs",
	"example",
	"examples",
	"maven-test",
	"node_modules",
	"scripts",
	"test",
	"testData",
	"test_scripts",
	"testdata",
	"testing",
	"testresources",
	"tests",
	"vendor",
}

func getDeps(fullPath string, ignoreDirs []string) ([]Dependency, Bitmask, error) {
	var discovered Discovered
	// special var so we don't double handle both repos with both
	// a Gemfile and Gemfile.lock
	var seenGemfile string

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, 0, os.ErrNotExist
	}

	pomPath := filepath.Join(fullPath, "pom.xml")
	// goPath := filepath.Join(fullPath, "go.mod")
	goPkgPath := filepath.Join(fullPath, "Gopkg.lock")
	glidePath := filepath.Join(fullPath, "glide.lock")
	rubyPath := filepath.Join(fullPath, "Gemfile") // Later we translate Gemfile.lock -> Gemfile to handle both cases
	pythonPath := filepath.Join(fullPath, "requirements.txt")

	ignoreDirs = append(ignoreDirs, defaultIgnore...)
	log.Debugf("directories ignored: %s", ignoreDirs)

	// point at the parent repo, but can't assume where the indicators will be
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// prevent panic by handling failure https://golang.org/pkg/path/filepath/#Walk
			return err
		}

		if info.IsDir() {
			// prevent walking down the vendors, docs, etc
			if slices.Contains(ignoreDirs, info.Name()) {
				log.Debugf("Skipping '%s', directory  name '%s' in ignore list", path, info.Name())
				return filepath.SkipDir
			}
		} else {
			// Two checks, one for filenames and the second switch for full
			// paths. Useful if we're looking for top of repo

			// comparisons here are made against the filename only, not full path
			// so matches will be found at any level of the file tree, not just top-level
			filename := info.Name()
			switch filename {
			case "go.mod":
				pkgs, err := scan.GetGolangDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					discovered.foundTypes.DepFoundAddFlag(LangGolang)
				}

				for path, goPkg := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    path,
						Files:   goPkg.Gofiles,
						Version: goPkg.Version,
					}
					discovered.deps = append(discovered.deps, d)
				}
			case "package-lock.json":
				// if theres not a yarn.lock fall thru
				if _, err := os.Stat(
					filepath.Join(
						filepath.Dir(path),
						"yarn.lock")); err == nil {
					return nil
				}
				fallthrough

			case "yarn.lock":
				pkgs, err := scan.GetNodeJSDeps(path)
				if err != nil {
					// ignore error
					log.Debugf("failed to scan for nodejs: %s", path)
					return nil
				}

				if len(pkgs) > 0 {
					discovered.foundTypes.DepFoundAddFlag(LangNodeJS)
				}

				for _, p := range pkgs {
					discovered.deps = append(discovered.deps,
						Dependency{
							DepType: LangNodeJS,
							Path:    p.Name,
							Version: p.Version,
							Files:   []string{},
						})
				}
			case "Cargo.lock":
				pkgs, err := scan.GetCrates(path)
				if err != nil {
					// ignore error
					log.Debugf("failed to scan rust crates: %s", path)
					return nil
				}

				discovered = addPackagesToDeps(discovered, pkgs, LangRust)
			default:
				ext := filepath.Ext(filename)
				// java
				switch ext {
				case ".zip":
					// be more aggressive with zip files, must contain something java ish
					if ok, _ := utils.ZipContainsJava(path); !ok {
						return nil
					}
					fallthrough
				case ".jar":
					fallthrough
				case ".war":
					fallthrough
				case ".ear":
					fallthrough
				case ".adm":
					fallthrough
				case ".hpi":
					file := strings.Replace(filepath.Base(path), ext, "", 1) // get filename, check if we can ignore
					if strings.HasSuffix(file, "-sources") || strings.HasSuffix(file, "-javadoc") {
						return nil
					}

					dir := filepath.Dir(path)
					pkgs, err := scan.GetJavaDeps(dir, ignoreDirs)
					if err == nil {

						if len(pkgs) > 0 {
							discovered.foundTypes.DepFoundAddFlag(LangJava)
						}

						for name, version := range pkgs {
							// just in case we report the full path to the dep
							name = strings.Replace(name, fullPath, "", 1)

							// if the dep ends with -javadoc or -sources, not really interested
							if !strings.HasSuffix(version, "-javadoc") && !strings.HasSuffix(version, "-sources") {
								discovered.deps = append(discovered.deps,
									Dependency{
										DepType: LangJava,
										Path:    name,
										Version: version,
										Files:   []string{},
									})
							}
						}
					}
				}
			}

			// translate Gemfile.lock -> Gemfile, to handle either case
			// but also avoid double-handling, i.e. scanning once for each file
			path = strings.Replace(path, "Gemfile.lock", "Gemfile", 1)

			// comparisons here are against the full filepath, so will not match if
			// these filenames are found in subdirectories, only the top level
			switch path {
			case goPkgPath:
				pkgs, err := scan.GetGoPkgDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					discovered.foundTypes.DepFoundAddFlag(LangGolang)
				}
				for _, goPkg := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    goPkg.Name,
						Version: goPkg.Version,
					}
					discovered.deps = append(discovered.deps, d)
				}
			case glidePath:
				pkgs, err := scan.GetGlideDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					discovered.foundTypes.DepFoundAddFlag(LangGolang)
				}
				for _, goPkg := range pkgs {
					d := Dependency{
						DepType: LangGolang,
						Path:    goPkg.Name,
						Version: goPkg.Version,
					}
					discovered.deps = append(discovered.deps, d)
				}
			case pomPath:
				dir := filepath.Dir(path)
				pkgs, err := scan.GetJavaDeps(dir, ignoreDirs)
				if err != nil {
					return err
				}

				discovered = addPackagesToDeps(discovered, pkgs, LangJava)
			case rubyPath:
				// To prevent double handling of both Gemfile and Gemfile.lock
				// Earier we translate Gemfile.lock -> Gemfile
				if path == seenGemfile {
					break
				}

				pkgs, err := scan.GetRubyDeps(path)
				if err != nil {
					return err
				}

				discovered = addPackagesToDeps(discovered, pkgs, LangRuby)
				seenGemfile = path
			case pythonPath:
				pkgs, err := scan.GetPythonDeps(path)
				if err != nil {
					return err
				}

				discovered = addPackagesToDeps(discovered, pkgs, LangPython)
			}

		}
		return nil
	})
	if err != nil {
		return nil, 0, err // should't matter
	}

	return discovered.deps, discovered.foundTypes, nil
}

// findBaseDir walks a directory tree through empty subdirs til it finds a directory with content
func findBaseDir(fullPath string) (string, error) {
	log.Debugf("Checking %s", fullPath)
	files, err := os.ReadDir(fullPath)
	if err != nil {
		return "", fmt.Errorf("Could not read: %s", err)
	}
	if len(files) == 1 && files[0].IsDir() {
		return findBaseDir(filepath.Join(fullPath, files[0].Name()))
	}
	return fullPath, nil
}

// GetDeps scans a given repository and returns all dependencies found in a DependencyList struct.
func GetDeps(fullPath string, ignoreDirs ...string) ([]Dependency, Bitmask, error) {
	fullPath, err := findBaseDir(fullPath)
	if err != nil {
		return nil, 0, err
	}

	deps, foundTypes, err := getDeps(fullPath, ignoreDirs)
	if err != nil {
		return deps, foundTypes, err
	}
	// if no deps found, check one level lower in 'src' directory
	// but ignore any new errors
	if len(deps) == 0 {
		fullPath = filepath.Join(fullPath, "src")
		if _, err := os.Stat(fullPath); err != nil {
			log.Debugf("No deps found, trying %s", fullPath)
			deps, foundTypes, _ = getDeps(fullPath, ignoreDirs)
		}
	}

	// de-duplicate
	unique := removeDuplicates(deps)

	return unique, foundTypes, err
}

func removeDuplicates(deps []Dependency) []Dependency {
	seen := map[string]bool{}
	filtered := []Dependency{}
	for _, dep := range deps {
		key := dep.ToString()
		if _, ok := seen[key]; !ok {
			seen[key] = true
			filtered = append(filtered, dep)
		}
	}
	return filtered
}
