FROM debian:bookworm-slim AS base

FROM golang:1.25.7-bookworm AS builder
COPY --from=base / /
WORKDIR /app
ADD . /app
RUN make build

FROM base

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/dist/aegisbot .
COPY --from=builder /app/_rules ./_rules

ENTRYPOINT ["./aegisbot"]
