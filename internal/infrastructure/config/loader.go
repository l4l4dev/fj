package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	applicationconfig "github.com/l4l4dev/fj/internal/application/config"
)

const configurationRelativePath = "fj/config.toml"

type environmentLookup func(string) (string, bool)

type tomlConfiguration struct {
	Instances []tomlInstance `toml:"instances"`
}

type tomlInstance struct {
	Name       string `toml:"name"`
	Endpoint   string `toml:"endpoint"`
	Credential string `toml:"credential"`
}

func Load() (applicationconfig.Configuration, error) {
	return load(os.LookupEnv, os.ReadFile)
}

func load(getenv environmentLookup, readFile func(string) ([]byte, error)) (applicationconfig.Configuration, error) {
	path, err := configurationPath(getenv)
	if err != nil {
		return applicationconfig.Configuration{}, err
	}

	contents, err := readFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return applicationconfig.Configuration{}, fmt.Errorf("configuration file not found: %s", path)
		}
		return applicationconfig.Configuration{}, fmt.Errorf("read configuration file: %w", err)
	}

	var decoded tomlConfiguration
	if _, err := toml.Decode(string(contents), &decoded); err != nil {
		return applicationconfig.Configuration{}, fmt.Errorf("configuration file contains malformed TOML")
	}

	configuration := applicationconfig.Configuration{Instances: make([]applicationconfig.Instance, len(decoded.Instances))}
	for index, instance := range decoded.Instances {
		configuration.Instances[index] = applicationconfig.Instance{
			Name:       instance.Name,
			Endpoint:   applicationconfig.Endpoint(instance.Endpoint),
			Credential: applicationconfig.CredentialReference(instance.Credential),
		}
	}
	if err := configuration.Validate(); err != nil {
		return applicationconfig.Configuration{}, fmt.Errorf("configuration is invalid: %w", err)
	}

	return configuration, nil
}

func configurationPath(getenv environmentLookup) (string, error) {
	base, xdgSet := getenv("XDG_CONFIG_HOME")
	if xdgSet && strings.TrimSpace(base) != "" {
		if !filepath.IsAbs(base) {
			return "", fmt.Errorf("XDG_CONFIG_HOME must be an absolute path")
		}
		return filepath.Join(base, configurationRelativePath), nil
	}

	home, homeSet := getenv("HOME")
	if !homeSet || strings.TrimSpace(home) == "" {
		return "", fmt.Errorf("HOME must be set when XDG_CONFIG_HOME is unset")
	}
	return filepath.Join(home, ".config", configurationRelativePath), nil
}
