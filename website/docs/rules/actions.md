---
title: Actions
description: The actions that can be taken and their parameters
---

When creating a rule you can give the following actions:

- `alert`
- `ban`
- `delete`
- `kick`
- `timeout`

## Parameters

Based on the action that will be taken, you can pass in the following additional parameters

### `alert`

| Parameter   | Description                                                        | Default                                                           |
| ----------- | ------------------------------------------------------------------ | ----------------------------------------------------------------- |
| `channelId` | The channel ID where the message will be sent                      | _Channel ID where the trigger occured_                            |
| `message`   | The message that will be sent, supports [templating](./templating) | `⚠️ Rule '{{.RuleName}}' matched and triggered by <@{{.UserID}}>` |

### `timeout`

| Parameter  | Description                                            | Default |
| ---------- | ------------------------------------------------------ | ------- |
| `duration` | The duration of the timeout, for example `15m` or `6h` | `10m`   |

### `ban`

| Parameter | Description                                                        | Default                                     |
| --------- | ------------------------------------------------------------------ | ------------------------------------------- |
| `reason`  | The reason that will be given, supports [templating](./templating) | `Matched rule '{{.RuleName}}', action=ban` |

### `kick`

| Parameter | Description                                                        | Default                                     |
| --------- | ------------------------------------------------------------------ | ------------------------------------------- |
| `reason`  | The reason that will be given, supports [templating](./templating) | `Matched rule '{{.RuleName}}', action=kick` |
