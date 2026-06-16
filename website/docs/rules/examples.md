---
title: Examples
description: Some rule examples to get inspired
---

## New accounts

```yaml
- rule:
    name: "NewAccounts"
    meta:
      event: "member_join"
      ignoreBots: true
    actions:
      - type: "alert"
        channelId: "1215772005681332395"
    expression: |
      time.since(member.created_at) < duration("168h")
```

## Crypto scam pictures

```yaml
- rule:
    name: "CryptoScamPictures"
    meta:
      event: "message_create"
      ignoreBots: true
    actions:
      - type: "timeout"
        duration: "6h"
    expression: |
      message.getLinks().size() == 4 &&
      message.getLinks().all(l,
        (
          l.startsWith("https://cdn.discordapp.com/attachments/") ||
          l.startsWith("https://media.discordapp.net/attachments/")
        ) &&
        l.matches(".*\\.(png|jpg|jpeg|gif|webp)(\\?.*)?$")
      )
```
