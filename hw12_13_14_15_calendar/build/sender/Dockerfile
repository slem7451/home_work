FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ./bin/sender ./cmd/calendar_sender

FROM alpine:3.9

COPY --from=build /app/bin/sender .
COPY /configs/sender_config.toml ./configs/sender_config.toml

CMD ["./sender", "-config", "./configs/sender_config.toml"]
