# RNOKPP Telegram Bot

🤖 RNOKPP Telegram Bot. It decodes an РНОКПП (Ukrainian tax number) into gender and date of birth, and generates random valid numbers.

The bot runs as an AWS Lambda behind a Function URL, and receives updates through a Telegram webhook.

[![Build Status](https://img.shields.io/github/actions/workflow/status/fre5h/rnokpp-telegram-bot/ci.yml?branch=main&style=flat-square)](https://github.com/fre5h/rnokpp-telegram-bot/actions?query=workflow%3ACI+branch%3Amain+)
[![CodeCov](https://img.shields.io/codecov/c/github/fre5h/rnokpp-telegram-bot.svg?style=flat-square)](https://codecov.io/github/fre5h/rnokpp-telegram-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/fre5h/rnokpp-telegram-bot?style=flat-square)](https://goreportcard.com/report/github.com/fre5h/rnokpp-telegram-bot)
[![License](https://img.shields.io/github/license/fre5h/rnokpp-telegram-bot?style=flat-square)](https://pkg.go.dev/github.com/fre5h/rnokpp-telegram-bot)

## Commands

- `/start` – help text.
- `/getinfo` – asks for an РНОКПП and replies with its gender and birthdate. Typing a number without a command does the same.
- `/random` – one random valid РНОКПП.
- `/randomn` – choose 2, 3, 5, 10 or 20 random valid РНОКПП.

## Configuration

| Variable             | Required | Description                                                         |
|----------------------|----------|---------------------------------------------------------------------|
| `TELEGRAM_BOT_TOKEN` | yes      | Bot token from @BotFather. The program refuses to start without it. |

## Development

```sh
go test -race ./...
```

## Build

```sh
make aws_lambda
```
