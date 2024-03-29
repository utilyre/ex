FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./ ./
RUN ./scripts/build.sh

FROM alpine:3.19

COPY --from=builder /app/build /app

ENV MODE="PROD"
ENV LOG_LEVEL="INFO"
ENV APP_ROOT="/app"
ENV SERVER_ADDR="0.0.0.0:80"

EXPOSE 80
CMD /app/server
