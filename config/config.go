package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rules struct {
		Disabled []string `yaml:"disabled"`
	} `yaml:"rules"`
}

var cfg Config

func LoadConfig() {
	data, err := os.ReadFile(".vibe-check.yaml")
	if err == nil {
		yaml.Unmarshal(data, &cfg)
	}
}

func IsRuleEnabled(rule string) bool {
	for _, disabled := range cfg.Rules.Disabled {
		if disabled == rule {
			return false
		}
	}
	return true
}
