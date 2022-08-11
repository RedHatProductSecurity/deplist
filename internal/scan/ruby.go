package scan

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// var RUBY_VERSIONS []string = []string{"system", "2.6.8"}

// GetRubyDeps uses `bundle update --bundler` to list ruby dependencies when a
// Gemfile.lock file exists
func GetRubyDeps(path string) (map[string]string, error) {
	return GetRubyDepsWithVersion(path, "system")
}

func setRbenvVersion(version string, cmd *exec.Cmd) {
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "RBENV_VERSION="+version)
}

func GetRubyDepsWithVersion(path, version string) (map[string]string, error) {
	if version != "system" {
		log.Debugf("retrying with ruby%v \n -> GetRubyDeps %s", version, path)
	} else {
		log.Debugf("GetRubyDeps %s", path)
	}

	gathered := make(map[string]string)

	dirPath := filepath.Dir(path)

	//Make sure that the Gemfile we are loading is supported by the version of bundle currently installed.
	cmd := exec.Command("bundle", "update", "--bundler")
	cmd.Dir = dirPath
	setRbenvVersion(version, cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Debug(string(output))
		if version == "2.6.8" {
			return nil, err
		}
		return GetRubyDepsWithVersion(path, "2.6.8")
	}

	cmd = exec.Command("bundle", "list")

	cmd.Dir = dirPath
	setRbenvVersion(version, cmd)

	data, err := cmd.Output()
	if err != nil {
		log.Debug(err)
		log.Debug(string(data))

		if version == "2.6.8" {
			return nil, err
		}
		return GetRubyDepsWithVersion(path, "2.6.8")
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
