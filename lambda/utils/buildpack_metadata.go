package utils

import (
	"fmt"
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/pelletier/go-toml"
)

type Configuration struct {
	Name        string `toml:"name"`
	Default     string `toml:"default"`
	Description string `toml:"description"`
}

type Metadata struct {
	Configurations []Configuration `toml:"configurations"`
}

type BuildpackMetadata struct {
	Buildpack packit.BuildpackInfo `toml:"buildpack"`
	Metadata  Metadata             `toml:"metadata"`
}

func ReadBuildpackMetadata(path string) (BuildpackMetadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return BuildpackMetadata{}, fmt.Errorf("unable to open buildpack.toml: %w", err)
	}
	defer file.Close()

	var buildpackMetadata BuildpackMetadata
	tomlParser := toml.NewDecoder(file)
	if err := tomlParser.Decode(&buildpackMetadata); err != nil {
		return BuildpackMetadata{}, fmt.Errorf("unable to decode buildpack.toml: %w", err)
	}

	return buildpackMetadata, nil
}
