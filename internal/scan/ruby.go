package scan

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//go:embed gemfile-parser.rb
var gemfileScriptData []byte

//go:embed gemspec-parser.rb
var gemSpecScriptData []byte

type Script struct {
	Name string
	Data []byte
}

var gemFileParser = Script{
	Name: "gemfile-parser.rb",
	Data: gemfileScriptData,
}

var gemSpecParser = Script{
	Name: "gemspec-parser.rb",
	Data: gemSpecScriptData,
}

func GetRubyDeps(path string) (map[string]string, error) {
	log.Debugf("GetRubyDeps %s", path)
	baseDir := filepath.Dir(path)
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}
	var gemspec string
	for _, e := range entries {
		filename := e.Name()
		if strings.HasSuffix(filename, ".gemspec") {
			gemspec = filepath.Join(baseDir, filename)
			break
		}
	}

	if gemspec != "" {
		log.Debugf("Found %s, parsing", gemspec)
		return runRubyParser(gemSpecParser, gemspec)
	}

	lockPath := filepath.Join(baseDir, "Gemfile.lock")

	if _, err := os.Stat(lockPath); err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Creating %s with `bundle lock`", lockPath)
			// Create Gemfile.lock
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, "env", fmt.Sprintf("--chdir=%s", baseDir), "bundle", "lock")
			data, err := cmd.CombinedOutput()
			if err != nil {
				log.Errorf("couldn't create %s: %v", lockPath, err)
				log.Debugf("bundle lock output: %v", string(data))
			}
			log.Debugf("Created %s", lockPath)
		} else {
			log.Errorf("Unexpected error: %v", err)
		}
	}
	return runRubyParser(gemFileParser, baseDir)
}

func runRubyParser(script Script, target string) (map[string]string, error) {
	gathered := make(map[string]string)

	g, err := os.CreateTemp("", script.Name)
	if err != nil {
		log.Errorf("Could not create ruby script %s: %s", script.Name, err)
		return gathered, err
	}
	defer os.Remove(g.Name())
	err = os.WriteFile(g.Name(), script.Data, 0644)
	if err != nil {
		log.Errorf("Could not write ruby script to %s: %s", g.Name(), err)
		return gathered, err
	}
	dir := filepath.Dir(target)
	name := filepath.Base(target)
	args := []string{fmt.Sprintf("--chdir=%s", dir), "ruby", g.Name(), name}
	log.Debugf("Running env %v", args)
	cmd := exec.Command("env", args...)
	data, err := cmd.Output()
	if err != nil {
		log.Errorf("Error running %s: %v: %s", script.Name, err, string(data))
		return gathered, err
	}

	splitOutput := strings.Split(string(data), "\n")

	for _, line := range splitOutput {
		if line == "" {
			continue
		}
		dep := strings.Split(line, " ")
		if len(dep) == 1 && dep[0] != "" {
			// no version found for this dep
			dep = append(dep, "")
		} else if len(dep) < 2 {
			log.Debugf("Unexpected dependency: %v, skipping", dep)
			continue
		}
		gathered[dep[0]] = dep[1]
	}

	return gathered, nil
}
