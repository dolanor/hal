package main

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type wakatimeConfig struct {
	APIKey string `toml:settings.api_key`
}

func loadWakatimeAPIKey() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	wakatimeCfgPath := filepath.Join(homePath, ".wakatime.cfg")

	iniFile, err := ini.LoadSources(ini.LoadOptions{
		AllowPythonMultilineValues: true,
	}, wakatimeCfgPath)
	if err != nil {
		return "", err
	}

	section, err := iniFile.GetSection("settings")
	if err != nil {
		return "", err
	}
	apiKey, err := section.GetKey("api_key")
	if err != nil {
		return "", err
	}

	return apiKey.String(), nil
}
