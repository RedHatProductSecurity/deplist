package scan

import (
	"embed"
	"os"
	"os/exec"
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

	if _, err := os.Stat("Gemfile.lock"); err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Creating Gemfile.lock with `bundle lock`")
			// Create Gemfile.lock
			cmd := exec.Command("bundle", "lock")
			data, err := cmd.CombinedOutput()
			if err != nil {
				log.Errorf("couldn't create Gemfile.lock: %v: %v", err, string(data))
				return nil, err
			}
			log.Debugf("Created Gemfile.lock")
		} else {
			log.Errorf("Unexpected error: %v", err)
			return nil, err
		}
	}
	return runGemlockParser()
}

func runGemlockParser() (map[string]string, error) {
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
	cmd := exec.Command("ruby", g.Name())
	data, err := cmd.Output()
	if err != nil {
		log.Errorf("Error running Gemfile.lock parser:  %v", err)
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
