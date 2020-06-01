package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	ini "gopkg.in/ini.v1"
	yaml "gopkg.in/yaml.v2"

	"github.com/boxofimagination/bxdk/go/env"
)

var (
	// ErrNoFileFound returned when there is no config file found
	ErrNoFileFound = errors.New("no config file found")
)


// Read configuration from the given paths.
// it will use the first file found in the given paths.
// It returns ErrNoFileFound if none of the given paths is exist
func Read(dest interface{}, paths ...string) error {
	for _, path := range paths {
		path = replacePathByEnv(path)

		// check if this path is exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		// load config
		ext := filepath.Ext(path)
		if ext == ".ini" {
			return loadIniConfig(dest, path)
		}
		return loadConfig(dest, path, ext)
	}
	return ErrNoFileFound
}

func loadIniConfig(dest interface{}, path string) error {
	f, err := ini.Load(path)
	if err != nil {
		return err
	}
	return f.MapTo(dest)
}

func loadConfig(dest interface{}, path, ext string) error {
	cfg := config.NewConfig()
	err := cfg.Load(file.NewSource(
		file.WithPath(path),
	))
	if err != nil {
		return err
	}

	// somehow yaml not working using micro-config
	// need more investigation regarding this issue
	if isYaml(ext) {
		return yaml.Unmarshal(cfg.Bytes(), dest)
	}

	return cfg.Scan(dest)
}

func replacePathByEnv(path string) string {
	boxEnv := string(env.ServiceEnv())
	return strings.Replace(path, "{BOXENV}", boxEnv, -1)
}

// isYaml returns true if the given extension is YAML file
func isYaml(ext string) bool  {
	return ext == ".yml" || ext == ".yaml"
}