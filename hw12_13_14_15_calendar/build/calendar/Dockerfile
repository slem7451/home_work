FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ./bin/calendar ./cmd/calendar

FROM alpine:3.9

COPY --from=build /app/bin/calendar .
COPY /configs/calendar_config.toml ./configs/calendar_config.toml

CMD ["./calendar", "-config", "./configs/calendar_config.toml"]
