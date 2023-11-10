package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// scheduleRe is a regexp for time range 00:00 - 23:59.
var scheduleRe = regexp.MustCompile("([0-1][0-9]|2[0-3]):[0-5][0-9]")

// Config represents application configuration.
type Config struct {
	Debug    bool           `yaml:"debug"`
	Schedule []Schedule     `yaml:"schedule"`
	Telegram TelegramConfig `yaml:"telegram"`
	ChatGPT  ChatGPTConfig  `yaml:"chatgpt"`
}

// Schedule describes what and when to run.
type Schedule struct {
	Time   string `yaml:"time"`
	Prompt string `yaml:"prompt"`
}

// TelegramConfig is a set of parameters for using Telegram API.
type TelegramConfig struct {
	APIKey  string `yaml:"api_key"`
	Channel string `yaml:"channel"`
}

// ChatGPTConfig is a set of parameters for using ChatGPT API.
type ChatGPTConfig struct {
	APIKey  string        `yaml:"api_key"`
	Timeout time.Duration `yaml:"timeout"`
	Prompts []Prompt      `yaml:"prompts"`
}

// Prompt describes a prompt for ChatGPT.
type Prompt struct {
	Name      string   `yaml:"name"`
	Template  string   `yaml:"template"`
	Arguments []string `yaml:"arguments"`
}

// String renders the prompt to a single string.
func (p *Prompt) String() string {
	arg := p.Arguments[rand.Intn(len(p.Arguments))]
	return fmt.Sprintf(p.Template, arg)
}

// ReadConfig returns configuration populated from the config file.
func ReadConfig(file string) (Config, error) {
	data, err := os.ReadFile(file) //nolint:gosec
	if err != nil {
		return Config{}, fmt.Errorf("read file: %w", err)
	}
	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal file: %w", err)
	}

	// Validate
	if len(conf.Schedule) == 0 {
		return Config{}, errors.New("empty schedule")
	}
	for _, s := range conf.Schedule {
		if !scheduleRe.MatchString(s.Time) {
			return Config{}, fmt.Errorf("invalid schedule: %s", s.Time)
		}
		if s.Prompt == "" {
			return Config{}, fmt.Errorf("empty prompt at %s", s.Time)
		}
	}
	if conf.ChatGPT.Timeout == 0 {
		conf.ChatGPT.Timeout = 1 * time.Minute
	}
	if len(conf.ChatGPT.Prompts) == 0 {
		return Config{}, errors.New("no prompts")
	}
	for _, p := range conf.ChatGPT.Prompts {
		if strings.Count(p.Template, "%s") != 1 {
			return Config{}, fmt.Errorf("invalid template %s", p.Name)
		}
	}

	return conf, nil
}
