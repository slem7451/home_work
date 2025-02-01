FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ./bin/scheduler ./cmd/calendar_scheduler

FROM alpine:3.9

COPY --from=build /app/bin/scheduler .
COPY /configs/scheduler_config.toml ./configs/scheduler_config.toml

CMD ["./scheduler", "-config", "./configs/scheduler_config.toml"]
