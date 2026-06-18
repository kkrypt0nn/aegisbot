---
title: Expression
description: The expression definition and their objects
---

The expression of the rule **must return a boolean** to know whether the rule gets triggered and the action taken.

The expression uses [CEL (Common Expression Language)](https://github.com/cel-expr/cel-spec) to be as dynamic and personalized as possible.

## Objects per event

The following objects can be used in the expression of the event

### `message_create` and `message_update`

- [`author`](/docs/objects/member) - Same as `message.author`
- [`message`](/docs/objects/message)

### `member_join` and `member_update`

- [`member`](/docs/objects/member)
