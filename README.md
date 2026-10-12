# rnokpp-telegram-bot

🤖 RNOKPP Telegram Bot. It decodes an РНОКПП (Ukrainian tax number) into gender and date of birth, and generates random valid numbers.

The bot runs as an AWS Lambda behind a Function URL, and receives updates through a Telegram webhook.

## Commands

- `/start` – help text.
- `/getinfo` – asks for an РНОКПП and replies with its gender and birth date. Typing a number without a command does the same.
- `/random` – one random valid РНОКПП.
- `/randomn` – choose 2, 3, 5, 10 or 20 random valid РНОКПП.

## Configuration

| Variable | Required | Description |
|---|---|---|
| `TELEGRAM_BOT_TOKEN` | yes | Bot token from @BotFather. The program refuses to start without it. |
| `TELEGRAM_WEBHOOK_SECRET` | recommended | Secret checked against the `X-Telegram-Bot-Api-Secret-Token` header. If unset, webhook requests are not authenticated. |

## Development

```sh
go test -race ./...
```

## Build and deploy

```sh
make aws_lambda
```

This builds a `bootstrap` binary for `linux/arm64` with the `provided` runtime and packs it into `archive.zip`. Upload the archive to a Lambda function with the `provided.al2023` (or `provided.al2`) runtime, architecture `arm64`, and enable a Function URL.
