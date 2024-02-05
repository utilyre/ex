FROM golang:1.21.6-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN ./scripts/build.sh

FROM alpine:3.19

ARG APP_ROOT=/app
ARG SERVER_PORT=80

ENV LOG_LEVEL=INFO
ENV APP_ROOT=${APP_ROOT}
ENV SERVER_ADDR=0.0.0.0:${SERVER_PORT}
ENV DSN=/app/data.db

COPY --from=builder /app/build ${APP_ROOT}

EXPOSE ${SERVER_PORT}
CMD /app/server -mode prod
