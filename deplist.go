package deplist

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/RedHatProductSecurity/deplist/internal/scan"
	"github.com/RedHatProductSecurity/deplist/internal/utils"

	log "github.com/sirupsen/logrus"
)

// enums start at 1 to allow us to specify found languages 0 = nil
const (
	LangGolang = 1 << iota
	LangJava
	LangNodeJS
	LangPython
	LangRuby
)

func init() {
	// check for the library required binaries
	languages := map[string]string{
		"yarn":   "yarn",
		"npm":    "npm",
		"go":     "go",
		"mvn":    "maven",
		"bundle": "bundler gem",
	}

	for lang_bin, lang_name := range languages {
		if _, err := exec.LookPath(lang_bin); err != nil {
			log.Fatal(lang_name, " is required in PATH")
		}
	}
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

func getDeps(fullPath string) ([]Dependency, Bitmask, error) {
	var discovered Discovered
	// special var so we don't double handle both repos with both
	// a Gemfile and Gemfile.lock
	var seenGemfile string

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, 0, os.ErrNotExist
	}

	pomPath := filepath.Join(fullPath, "pom.xml")
	goPath := filepath.Join(fullPath, "go.mod")
	goPkgPath := filepath.Join(fullPath, "Gopkg.lock")
	glidePath := filepath.Join(fullPath, "glide.lock")
	rubyPath := filepath.Join(fullPath, "Gemfile") // Later we translate Gemfile.lock -> Gemfile to handle both cases
	pythonPath := filepath.Join(fullPath, "requirements.txt")

	// point at the parent repo, but can't assume where the indicators will be
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// prevent panic by handling failure https://golang.org/pkg/path/filepath/#Walk
			return err
		}

		if info.IsDir() {
			// prevent walking down the docs, .git, tests, etc.
			if utils.BelongsToIgnoreList(info.Name()) {
				return filepath.SkipDir
			}
		} else {
			// Two checks, one for filenames and the second switch for full
			// paths. Useful if we're looking for top of repo
			switch filename := info.Name(); filename {
			// for now only go for yarn and npm
			case "package.json":
				pkg, err := scan.GetNodeJSPackage(path)
				if err != nil {
					log.Debugf("failed to scan for nodejs package: %s", path)
					return nil
				}

				foundTypes.DepFoundAddFlag(LangNodeJS)

				deps = append(deps,
					Dependency{
						DepType:   LangNodeJS,
						Path:      pkg.Name,
						Version:   pkg.Version,
						Files:     []string{},
						IsBundled: true,
					})
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

					pkgs, err := scan.GetJarDeps(path)
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
										DepType:   LangJava,
										Path:      name,
										Version:   version,
										Files:     []string{},
										IsBundled: true,
									})
							}
						}
					}
				}
			}

			// translate Gemfile.lock -> Gemfile, to handle either case
			// but also avoid double-handling, i.e. scanning once for each file
			path = strings.Replace(path, "Gemfile.lock", "Gemfile", 1)

			switch path {
			case goPath: // just support the top level go.mod for now
				pkgs, err := scan.GetGolangDeps(path)
				if err != nil {
					return err
				}

				if len(pkgs) > 0 {
					discovered.foundTypes.DepFoundAddFlag(LangGolang)
				}

				for path, goPkg := range pkgs {
					d := Dependency{
						DepType:   LangGolang,
						Path:      path,
						Files:     goPkg.Gofiles,
						Version:   goPkg.Version,
						IsBundled: true,
					}
					discovered.deps = append(discovered.deps, d)
				}
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
						DepType:   LangGolang,
						Path:      goPkg.Name,
						Version:   goPkg.Version,
						IsBundled: true,
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
						DepType:   LangGolang,
						Path:      goPkg.Name,
						Version:   goPkg.Version,
						IsBundled: true,
					}
					discovered.deps = append(discovered.deps, d)
				}
			case pomPath:
				pkgs, err := scan.GetMvnDeps(path)
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
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return "", fmt.Errorf("Could not read: %s", err)
	}
	if len(files) == 1 && files[0].IsDir() {
		return findBaseDir(filepath.Join(fullPath, files[0].Name()))
	}
	return fullPath, nil
}

// GetDeps scans a given repository and returns all dependencies found in a DependencyList struct.
func GetDeps(fullPath string) ([]Dependency, Bitmask, error) {
	fullPath, err := findBaseDir(fullPath)
	if err != nil {
		return nil, 0, err
	}

	deps, foundTypes, err := getDeps(fullPath)
	if err != nil {
		return deps, foundTypes, err
	}
	// if no deps found, check one level lower in 'src' directory
	// but ignore any new errors
	if len(deps) == 0 {
		fullPath = filepath.Join(fullPath, "src")
		log.Debugf("Checking %s", fullPath)
		deps, foundTypes, _ = getDeps(fullPath)
	}

	return deps, foundTypes, err
}
