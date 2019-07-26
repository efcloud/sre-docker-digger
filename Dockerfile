ARG SOURCE=golang:1.11

# Lint stage
FROM $SOURCE as lint

RUN make in-docker-lint

# Test stage
FROM $SOURCE as test

RUN make in-docker-test

# Build the app
FROM $SOURCE as builder

RUN make in-docker-build-app

# Release container
FROM alpine:3.9 as final

ENV DATADOG_HOST=https://api.datadoghq.eu

WORKDIR /app

COPY --from=builder /go/src/digger/digger /usr/local/bin/digger


ENTRYPOINT ["/usr/local/bin/digger"]
CMD [""]
