package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	openai "github.com/sashabaranov/go-openai"
)

// Bot is the main app entity. It runs periodically to generate and
// publish content.
type Bot struct {
	schedule map[string]string
	prompts  map[string]Prompt
	debug    bool
	telegram *tg.BotAPI
	channel  string
	chatgpt  *openai.Client
	timeout  time.Duration
}

// NewBot creates new bot.
func NewBot(conf Config) (*Bot, error) {
	if !strings.HasPrefix(conf.Telegram.Channel, "@") {
		conf.Telegram.Channel = "@" + conf.Telegram.Channel
	}
	t, err := tg.NewBotAPI(conf.Telegram.APIKey)
	if err != nil {
		return nil, fmt.Errorf("init telegram API: %w", err)
	}

	ai := openai.NewClient(conf.ChatGPT.APIKey)

	schedule := map[string]string{}
	for _, s := range conf.Schedule {
		schedule[s.Time+":00"] = s.Prompt
	}
	prompts := map[string]Prompt{}
	for _, p := range conf.ChatGPT.Prompts {
		prompts[p.Name] = p
	}

	bot := Bot{
		schedule: schedule,
		prompts:  prompts,
		debug:    conf.Debug,
		telegram: t,
		channel:  conf.Telegram.Channel,
		chatgpt:  ai,
		timeout:  conf.ChatGPT.Timeout,
	}

	return &bot, nil
}

// Run starts periodical generate-publish runs.
func (b *Bot) Run(ctx context.Context) {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()

	for range timer.C {
		if ctx.Err() != nil {
			return
		}
		promptName, ok := b.timeToRun()
		if !ok {
			continue
		}
		if err := b.RunOnce(ctx, promptName); err != nil {
			log.Printf("Bot failed with error: %v", err)
			continue
		}
	}
}

// RunOnce generates text and publishes it one single time.
func (b *Bot) RunOnce(ctx context.Context, promptName string) error {
	p, ok := b.prompts[promptName]
	if !ok {
		return fmt.Errorf("unknown prompt: %s", promptName)
	}
	prompt := p.String()
	if b.debug {
		log.Printf("Prompt:\n%s", prompt)
	}
	text, err := b.generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("generate text: %w", err)
	}
	if b.debug {
		log.Printf("Text:\n%s", text)
		return nil
	}
	if err := b.publish(ctx, text); err != nil {
		return fmt.Errorf("publish text: %w", err)
	}
	return nil
}

func (b *Bot) timeToRun() (string, bool) {
	now := time.Now().Format("15:04:05")
	promptName, ok := b.schedule[now]
	return promptName, ok
}

// generate generates a text using ChatGPT API.
func (b *Bot) generate(ctx context.Context, prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, b.timeout)
	defer cancel()

	req := openai.ChatCompletionRequest{
		// Model: "gpt-4o-mini",
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		}},
	}

	resp, err := b.chatgpt.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("completion error: %w", err)
	}
	return resp.Choices[0].Message.Content, nil
}

// publish sends a text to Telegram.
func (b *Bot) publish(_ context.Context, text string) error {
	msg := tg.NewMessageToChannel(b.channel, text)
	_, err := b.telegram.Send(msg)
	return err //nolint:wrapcheck
}
