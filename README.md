# [wikisourcebot](https://t.me/Wikisource_bot)

[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](./LICENSE) [![License: GPL v3](https://img.shields.io/badge/Dev-%20@hearot-blue.svg)](https://t.me/hearot)

A Telegram bot which allows you to retrieve articles from wikisource.org.

To run the bot yourself, you will need:
- Go
- [mwclient](https://github.com/cgt/go-mwclient)
- [telebot](https://github.com/tucnak/telebot)

## Setup
- Get a token from [@BotFather](http://t.me/BotFather).
- Activate the *Inline mode* using the `/setinline` command with [@BotFather](http://t.me/BotFather).
- Set `token` in `bot.go`.
- Build the source code using `go build bot.go`.
- Finally, run the executable file you got from the previous step.
