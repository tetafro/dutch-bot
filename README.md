# dutch

[![Codecov](https://codecov.io/gh/tetafro/dutch-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/tetafro/dutch-bot)
[![Go Report](https://goreportcard.com/badge/github.com/tetafro/dutch-bot)](https://goreportcard.com/report/github.com/tetafro/dutch-bot)
[![CI](https://github.com/tetafro/dutch-bot/actions/workflows/push.yml/badge.svg)](https://github.com/tetafro/dutch-bot/actions)

A small bot for learning Dutch. It generates grammar rules explanations for a
random topic (using ChatGPT API). Texts are posted to a Telegram channel.

## Build and run

Create a bot and get Telegram API token from the bot called `@botfather` (free).

Get OpenAI API key [here](https://platform.openai.com/account/api-keys) (paid).

Copy and populate config
```sh
cp config.example.yaml config.yaml
```

Build
```sh
make build
```

Run in a loop on schedule, results are published to Telegram
```sh
./bin/dutch-bot
```

Generate text and print results without publishing to Telegram
```sh
./bin/dutch-bot -once -debug
```

## Encrypted config

Encrypt
```sh
echo "password" > .vault_pass.txt
ansible-vault encrypt \
    --output config.yaml.vault \
    config.yaml
```

Edit
```sh
EDITOR='code --wait' \
ansible-vault edit config.yaml.vault
```

## Deploy

Normally deploy is done by Github actions.

Manual deploy
```sh
SSH_SERVER=10.0.0.1:22 \
SSH_USER=user \
make deploy
```
