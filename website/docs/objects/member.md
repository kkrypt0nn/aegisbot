---
title: Member
description: The member object with its fields and methods
---

## Fields

| Field        | Type      | Description                                    |
| ------------ | --------- | ---------------------------------------------- |
| `username`   | String    | The username of the member.                    |
| `bot`        | Boolean   | Whether the member is a bot.                   |
| `created_at` | Timestamp | When the member created their Discord account. |
| `joined_at`  | Timestamp | When the member joined the Discord server.     |

## Methods

### `isBot()`

Returns a boolean whether the member is a bot.
