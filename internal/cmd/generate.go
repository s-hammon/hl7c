package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/s-hammon/hl7c/internal/config"
)

func readConfig(stderr io.Writer, dir, filename string) (string, *config.Config, error) {
	configPath := ""
	if filename != "" {
		configPath = filepath.Join(dir, filename)
	} else {
		yamlPath := filepath.Join(dir, "model_config.yaml")

		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			fmt.Fprintln(stderr, "error parsing model_config.yaml: file does not exit")
			return "", nil, errors.New("model_config.yaml file missing")
		}

		configPath = yamlPath
	}

	base := filepath.Base(configPath)
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: file does not exist\n", base)
		return "", nil, err
	}
	defer file.Close()

	conf, err := config.ParseConfig(file)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %s: %s\n", base, err)
		return "", nil, err
	}

	return conf.Meta.Package, &conf, nil
}

func Generate(ctx context.Context, dir, filename string, o *Options) (map[string]string, error) {
	pkgDir, conf, err := o.ReadConfig(dir, filename)
	if err != nil {
		return nil, err
	}

	models := make(map[string]string)
	modelDir := filepath.Join("internal/", pkgDir)
	models[modelDir] = conf.Compile(filepath.Base(pkgDir))
	return models, nil
}
