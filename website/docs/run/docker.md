---
title: Docker
description: Run Aegisbot with Docker
---

You can run the tool from the published [Docker image](https://ghcr.io/kkrypt0nn/aegisbot) using

```bash
docker run -it \
    -e BOT_TOKEN="..." \
    -e RULES_FOLDER="/app/_rules" \
    -v ./_rules:/app/_rules \
    ghcr.io/kkrypt0nn/aegisbot
```
