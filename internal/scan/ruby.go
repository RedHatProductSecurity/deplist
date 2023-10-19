package scan

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const scriptName = "gemlock-parser.rb"

// just here to satisfy gofmt
var assets embed.FS

//go:embed gemlock-parser.rb
var rubyScript []byte

func GetRubyDeps(path string) (map[string]string, error) {
	log.Debugf("GetRubyDeps %s", path)
	baseDir := filepath.Dir(path)
	lockPath := filepath.Join(baseDir, "Gemfile.lock")

	if _, err := os.Stat(lockPath); err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Creating %s with `bundle lock`", lockPath)
			// Create Gemfile.lock
			cmd := exec.Command("env", fmt.Sprintf("--chdir=%s", baseDir), "bundle", "lock")
			data, err := cmd.CombinedOutput()
			if err != nil {
				log.Errorf("couldn't create %s: %v: %v", lockPath, err, string(data))
				return nil, err
			}
			log.Debugf("Created %s", lockPath)
		} else {
			log.Errorf("Unexpected error: %v", err)
			return nil, err
		}
	}
	return runGemlockParser(lockPath)
}

func runGemlockParser(lockPath string) (map[string]string, error) {
	gathered := make(map[string]string)

	g, err := os.CreateTemp("", scriptName)
	if err != nil {
		log.Errorf("Could not create ruby script %s: %s", scriptName, err)
		return gathered, err
	}
	err = os.WriteFile(g.Name(), rubyScript, 0644)
	if err != nil {
		log.Errorf("Could not write ruby script to %s: %s", g.Name(), err)
		return gathered, err
	}
	args := []string{g.Name(), lockPath}
	log.Debugf("Running ruby %v", args)
	cmd := exec.Command("ruby", args...)
	data, err := cmd.Output()
	if err != nil {
		log.Errorf("Error running Gemfile.lock parser: %v: %s", err, string(data))
		return gathered, err
	}

	splitOutput := strings.Split(string(data), "\n")

	for _, line := range splitOutput {
		if line == "" {
			continue
		}
		dep := strings.Split(line, " ")
		if len(dep) < 2 {
			log.Debugf("Unexpected dependency: %v, skipping", dep)
			continue
		}
		gathered[dep[0]] = dep[1]
	}

	return gathered, nil
}
