package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		file, err := os.CreateTemp(os.TempDir(), "dutch-bot.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString(strings.Join([]string{
			`schedule:`,
			`  - time: "06:00"`,
			`    prompt: one`,
			`  - time: "11:00"`,
			`    prompt: one`,
			`  - time: "19:00"`,
			`    prompt: two`,
			`telegram:`,
			`  api_key: api-key`,
			`  channel: channel_name`,
			`chatgpt:`,
			`  api_key: api-key`,
			`  timeout: 1m`,
			`  prompts:`,
			`    - name: one`,
			`      template: "Template one: %s"`,
			`      arguments: [arg1]`,
			`    - name: two`,
			`      template: "Template two: %s"`,
			`      arguments: [arg2]`,
		}, "\n"))
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		conf, err := ReadConfig(file.Name())
		require.NoError(t, err)

		want := Config{
			Schedule: []Schedule{
				{Time: "06:00", Prompt: "one"},
				{Time: "11:00", Prompt: "one"},
				{Time: "19:00", Prompt: "two"},
			},
			Telegram: TelegramConfig{
				APIKey:  "api-key",
				Channel: "channel_name",
			},
			ChatGPT: ChatGPTConfig{
				APIKey:  "api-key",
				Timeout: 1 * time.Minute,
				Prompts: []Prompt{
					{
						Name:      "one",
						Template:  "Template one: %s",
						Arguments: []string{"arg1"},
					},
					{
						Name:      "two",
						Template:  "Template two: %s",
						Arguments: []string{"arg2"},
					},
				},
			},
		}
		require.Equal(t, want, conf)
	})

	t.Run("empty schedule parameter", func(t *testing.T) {
		file, err := os.CreateTemp(os.TempDir(), "dutch-bot.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString(strings.Join([]string{
			`schedule: []`,
			`telegram:`,
			`  api_key: api-key`,
			`  channel: channel_name`,
			`chatgpt:`,
			`  api_key: api-key`,
			`  prompts:`,
			`    - name: one`,
			`      template: "Template one: %s"`,
			`      arguments: [arg1]`,
		}, "\n"))
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		_, err = ReadConfig(file.Name())
		require.ErrorContains(t, err, "empty schedule")
	})

	t.Run("invalid schedule parameter", func(t *testing.T) {
		file, err := os.CreateTemp(os.TempDir(), "dutch-bot.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString(strings.Join([]string{
			`schedule:`,
			`  - time: "55:89"`,
			`    prompt: one`,
			`telegram:`,
			`  api_key: api-key`,
			`  channel: channel_name`,
			`chatgpt:`,
			`  api_key: api-key`,
			`  prompts:`,
			`    - name: one`,
			`      template: "Template one: %s"`,
			`      arguments: [arg1]`,
		}, "\n"))
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		_, err = ReadConfig(file.Name())
		require.ErrorContains(t, err, "invalid schedule")
	})

	t.Run("missing placeholder in prompt", func(t *testing.T) {
		file, err := os.CreateTemp(os.TempDir(), "dutch-bot.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString(strings.Join([]string{
			`schedule:`,
			`  - time: "11:00"`,
			`    prompt: one`,
			`telegram:`,
			`  api_key: api-key`,
			`  channel: channel_name`,
			`chatgpt:`,
			`  api_key: api-key`,
			`  prompts:`,
			`    - name: one`,
			`      template: "Template one"`,
			`      arguments: [arg1]`,
		}, "\n"))
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		_, err = ReadConfig(file.Name())
		require.ErrorContains(t, err, "invalid template")
	})

	t.Run("too many placeholders in prompt", func(t *testing.T) {
		file, err := os.CreateTemp(os.TempDir(), "dutch-bot.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString(strings.Join([]string{
			`schedule:`,
			`  - time: "11:00"`,
			`    prompt: one`,
			`telegram:`,
			`  api_key: api-key`,
			`  channel: channel_name`,
			`chatgpt:`,
			`  api_key: api-key`,
			`  prompts:`,
			`    - name: one`,
			`      template: "Template one %s %s"`,
			`      arguments: [arg1]`,
		}, "\n"))
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		_, err = ReadConfig(file.Name())
		require.ErrorContains(t, err, "invalid template")
	})

	t.Run("invalid yaml", func(t *testing.T) {
		file, err := os.CreateTemp(os.TempDir(), "dutch-bot.yaml")
		require.NoError(t, err)
		defer os.Remove(file.Name())

		_, err = file.WriteString("hello: world: !")
		require.NoError(t, err)

		err = file.Close()
		require.NoError(t, err)

		_, err = ReadConfig(file.Name())
		require.ErrorContains(t, err, "unmarshal file")
	})

	t.Run("non-existing config file", func(t *testing.T) {
		_, err := ReadConfig("not-exists.yml")
		require.ErrorContains(t, err, "no such file or directory")
	})
}
