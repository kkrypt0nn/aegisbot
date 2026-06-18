---
title: Message
description: The message object with its fields and methods
---

## Fields

| Field         | Type                                             | Description                                   |
| ------------- | ------------------------------------------------ | --------------------------------------------- |
| `content`     | String                                           | The content of the message.                   |
| `author`      | [`Member`](/docs/objects/member)                 | The member that is the author of the message. |
| `channel`     | [`Channel`](/docs/objects/channel)               | The channel in which the message is.          |
| `mentions`    | List of string                                   | The list of mentions in the message.          |
| `attachments` | List of [`Attachment`](/docs/objects/attachment) | The files that are attached to the message.   |

## Methods

### `isDM()`

Returns a boolean whether the message was sent in a `DM` channel.

### `getMentions()`

Returns the list of mentions.

### `hasLinks()`

Returns a boolean whether the message contains links.

### `getLinks()`

Returns the list of links.
