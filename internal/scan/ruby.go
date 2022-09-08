package scan

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

var RubyVersions []string = []string{"system"}

func GetInstalledRubyVersions() []string {
	cmd := exec.Command("rbenv", "versions", "--bare")
	data, err := cmd.Output()
	if err != nil {
		return nil
	}

	versions := strings.Split(string(data), "\n")
	versions = versions[:len(versions)-1]
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))

	RubyVersions = append(RubyVersions, versions...)

	if len(RubyVersions) == 1 {
		log.Debug("rbenv not detected, falling back to system ruby ONLY. Please ensure that bundler is installed and available in your path.")
	}

	return versions
}

func setRubyVersion(version string, cmd *exec.Cmd) {
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "RBENV_VERSION="+version)
}

// GetRubyDeps calls GetRubyDepsWithVersion with the system ruby version
func GetRubyDeps(path string) (map[string]string, error) {
	if len(RubyVersions) != 1 {
		log.Debugf("Ruby versions detected: %+v\n", RubyVersions)

		for _, version := range RubyVersions {
			cmd := exec.Command("gem", "install", "bundler")
			setRubyVersion(version, cmd)
			data, err := cmd.CombinedOutput()
			if err != nil {
				log.Debugf("couldn't install bundler: %v", string(data))
			}
			log.Debugf("Installed bundler for ruby %v\n", version)
		}
	}

	return GetRubyDepsWithVersion(path, 0)
}

// GetRubyDepsWithVersion uses `bundle list` to list ruby dependencies when a Gemfile.lock file exists
func GetRubyDepsWithVersion(path string, version int) (map[string]string, error) {
	if version >= len(RubyVersions) {
		log.Debug("GetRubyDeps Failed! No more ruby versions available")
		return nil, errors.New("GetRubyDeps Failed: all ruby versions failed " + path)
	}

	if version != 0 {
		log.Debug("retrying...")
	}

	log.Debugf("GetRubyDeps(%v) %s", RubyVersions[version], path)

	gathered := make(map[string]string)

	dirPath := filepath.Dir(path)

	// override the gem path otherwise might hit perm issues and it's annoying
	gemPath, err := os.MkdirTemp("", "gem_vendor")
	if err != nil {
		return nil, err
	}

	// cleanup after ourselves
	defer os.RemoveAll(gemPath)

	//Make sure that the Gemfile we are loading is supported by the version of bundle currently installed.
	cmd := exec.Command("bundle", "update", "--bundler")
	cmd.Dir = dirPath
	cmd.Env = append(cmd.Env, "BUNDLE_PATH="+gemPath)
	setRubyVersion(RubyVersions[version], cmd)

	data, err := cmd.CombinedOutput()
	if err != nil {
		log.Debugf("error: %s", err)
		log.Debugf("error: %v", string(data))
		return GetRubyDepsWithVersion(path, version+1)
	}

	cmd = exec.Command("bundle", "list")

	cmd.Dir = dirPath
	cmd.Env = append(cmd.Env, "BUNDLE_PATH="+gemPath)

	setRubyVersion(RubyVersions[version], cmd)

	data, err = cmd.Output()
	if err != nil {

		if version == len(RubyVersions) {
			log.Debugf("err: %v", err)
			log.Debugf("data: %v", string(data))
			return nil, err
		}
		return GetRubyDepsWithVersion(path, version+1)
	}

	splitOutput := strings.Split(string(data), "\n")

	for _, line := range splitOutput {
		if !strings.HasPrefix(line, "  *") {
			continue
		}
		rawDep := strings.TrimPrefix(line, "  * ")
		dep := strings.Split(rawDep, " ")
		dep[1] = dep[1][1 : len(dep[1])-1]
		gathered[dep[0]] = dep[1]
	}

	return gathered, nil
}
