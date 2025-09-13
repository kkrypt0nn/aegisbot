# Aegisbot

[![CI Badge](https://github.com/kkrypt0nn/aegisbot/actions/workflows/ci.yml/badge.svg)](https://github.com/kkrypt0nn/aegisbot/actions)
[![Discord Server Badge](https://img.shields.io/discord/1358456011316396295?logo=discord)](https://discord.gg/xj6y5ZaTMr)
[![Last Commit Badge](https://img.shields.io/github/last-commit/kkrypt0nn/aegisbot)](https://github.com/kkrypt0nn/aegisbot/commits/main)
[![Conventional Commits Badge](https://img.shields.io/badge/Conventional%20Commits-1.0.0-%23FE5196?logo=conventionalcommits&logoColor=white)](https://conventionalcommits.org/en/v1.0.0/)

> [!CAUTION]
> **This project is under active development and is not yet ready for production use.**
> I strongly advise **not** using it until a stable release is published.
> **Contributions (especially the creation of rules) and feedback are welcome and appreciated once the core is complete.**

### 🛡️ Pattern-match your Discord and Twitch defense

Aegisbot is a novel Discord and Twitch bot with **advanced pattern-matching auto-moderation**, built on concepts inspired by [YARA](https://virustotal.github.io/yara/) - not just yet another "auto-mod" clone.
It allows for fine-grained detection of malicious, spammy, or unwanted behavior using customizable matching rules written in a simple and common syntax such as YAML or JSON, rather than static keyword lists or simplistic triggers.

Huge thanks to [@evilsocket](https://github.com/evilsocket) for giving this bot idea.

## Getting Started

### Installation

> 🚧 Installation instructions will be added once the bot reaches a usable development milestone.

For now, feel free to watch the repository or join the [Discord server](https://discord.gg/xj6y5ZaTMr) to stay updated.

### Prerequisites

- Go ≥ 1.24
- A Discord bot token

### Run Locally

```bash
git clone https://github.com/kkrypt0nn/aegisbot
cd aegisbot
go run .
```

## Features

Coming soon!

## TODOs

Aegisbot will (hopefully) support:

- Pre-defined matches, like those weird and spammy characters
- Rule tags
- Scanning of messages, user profiles, rate of messages, etc.
- A GUI
  - To add, remove, deactivate and edit the rules
  - To run the bot super easily for local development

## Rule example

Currently a rule may look like

```yaml
- rule:
    name: "PhishingLink"
    meta:
      action: "alert"
      context: "message"
      ignoreBots: true
    strings:
      - name: "link"
        value: "https://badsite.com"
      - name: "scam"
        value: "free nitro"
    expression: |
      message.content.contains(link) && message.content.contains(scam)
```

## Documentation

Coming soon! Once the base rule engine and at least the Discord integrations are complete, a documentation will be hosted in the `docs/` folder and on a documentation website.

## Troubleshooting

Issues are currently **disabled** while Aegisbot is in active development and not ready for actual use.

Once a usable version is released, the issue tracker will be opened for bug reports, feature requests, and support.
In the meantime, feel free to join the [Discord Community](https://discord.gg/xj6y5ZaTMr) to follow development updates or ask questions.

## Contributing

Contributions are more than welcome, but please wait until we release a working beta or alpha version, then the issues will also be opened.

When it's time, please follow:
- [Contributing Guidelines](./CONTRIBUTING.md)
- [Code of Conduct](./CODE_OF_CONDUCT.md)

## License

This bot was made with 💜 by Krypton and is under the [AGPLv3 License](./LICENSE.md).
