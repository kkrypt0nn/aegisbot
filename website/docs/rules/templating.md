---
title: Templating
description: The templating opportunities for messages and descriptions of actions
---

When settin up some actions, for example `alert`, to be taken when a rule is triggered you will be able to set a custom text to be used. This text uses some [Go Templating](https://developer.hashicorp.com/nomad/docs/reference/go-template-syntax) capabilities and let you use some variables in the messages.

## Variables

Variables can be used with `{{.<variable>}}`, for example `{{.RuleName}}`. The available variables are:

| Variable    | Description                                |
| ----------- | ------------------------------------------ |
| `RuleName`  | The rule name that was triggered.          |
| `GuildID`   | The guild ID in which the message is in.   |
| `ChannelID` | The channel ID in which the message is in. |
| `MessageID` | The message ID that triggered the rule.    |
| `UserID`    | The user ID that triggered the rule.       |
