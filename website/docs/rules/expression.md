---
title: Expression
description: The expression definition and their objects
---

The expression of the rule **must return a boolean** to know whether the rule gets triggered and the action taken.

The expression uses [CEL (Common Expression Language)](https://github.com/cel-expr/cel-spec) to be as dynamic and personalized as possible.

## Objects per event

The following objects can be used in the expression of the event

### `message_create` and `message_update`

- [`author`](#member) - Same as `message.author`
- [`message`](#message)

### `member_join` and `member_update`

- [`member`](#member)

## Objects

The following objects with their fields and methods are defined

### `Attachment`

#### Fields

- `filename` - String
- `url` - String

### `Channel`

#### Fields

- `name` - String
- `type` - `ChannelType` enum

### `Member`

#### Fields

- `username` - String
- `bot` - Boolean
- `created_at` - Timestamp
- `joined_at` - Timestamp

#### Methods

- `isBot()` - Returns a boolean

### `Message`

#### Fields

- `content` - String
- `author` - [`Member`](#member)
- `channel` - [`Channel`](#channel)
- `mentions` - List of string
- `attachments` - List of [`Attachment`](#attachment)

#### Methods

- `isDM()` - Returns a boolean
- `getMentions()` - Returns a list of string
- `hasLinks()` - Returns a boolean
- `getLinks()` - Returns a list of string
