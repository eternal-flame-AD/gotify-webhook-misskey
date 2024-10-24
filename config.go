package main

import (
	"crypto/rand"
	"errors"
	"fmt"
)

const DummySecret = "ChangeMe"

type Config struct {
	Sources []SourceConfig `yaml:"sources"`
}

type SourceConfig struct {
	Slug    string `yaml:"slug"`
	Name    string `yaml:"name"`
	Secret  string `yaml:"secret"`
	BaseURL string `yaml:"base_url"`
}

func CreateDefaultConfig() Config {
	return Config{
		Sources: []SourceConfig{
			{
				Slug:    "me:me.com",
				Name:    "Notification for @me@me.com",
				Secret:  DummySecret,
				BaseURL: "https://me.invalid/",
			},
		},
	}
}

func (c *Config) GetSource(slug string) *SourceConfig {
	for i, source := range c.Sources {
		if source.Slug == slug {
			return &c.Sources[i]
		}
	}
	return nil
}

func randomSecret() string {
	var secret = make([]byte, 16)
	_, _ = rand.Read(secret)
	return fmt.Sprintf("%x", secret)
}

func allAlnum(s string) bool {
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func (c *Config) Validate() error {
	var errs []error
	for _, source := range c.Sources {
		if !allAlnum(source.Slug) {
			errs = append(errs, fmt.Errorf("source %s has invalid slug, only alnum characters are allowed", source.Slug))
		}
		if source.Secret == "" {
			errs = append(errs, fmt.Errorf("source %s has no secret, what about %s?", source.Slug, randomSecret()))
		}
		if source.Secret == DummySecret {
			errs = append(errs, fmt.Errorf("source %s has the default secret, please change it, what about %s?", source.Slug, randomSecret()))
		}
	}
	return errors.Join(errs...)
}
