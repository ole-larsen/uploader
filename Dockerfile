# syntax=docker/dockerfile:1
# STEP 1 dapp binary using golang:alpine-latest
FROM golang:alpine as go-builder

ARG APP_NAME=uploader
ARG PORT=1112
ARG NODE_ENV=production

# PostgreSQL settings (default)
ARG DB_SQL_HOST
ARG DB_SQL_PORT
ARG DB_SQL_USERNAME
ARG DB_SQL_PASSWORD
ARG DB_SQL_DATABASE

ARG USE_HASH=true
ARG USE_DB=true
ARG SESSION_SECRET
ARG X_TOKEN

# set environment variables
ENV APP_NAME=$APP_NAME
ENV PORT=$PORT
ENV NODE_ENV=$NODE_ENV

# PostgreSQL settings
ENV DB_SQL_HOST=$DB_SQL_HOST
ENV DB_SQL_PORT=$DB_SQL_PORT
ENV DB_SQL_USERNAME=$DB_SQL_USERNAME
ENV DB_SQL_PASSWORD=$DB_SQL_PASSWORD
ENV DB_SQL_DATABASE=$DB_SQL_DATABASE

ENV SESSION_SECRET=$SESSION_SECRET
ENV X_TOKEN=$X_TOKEN
ENV USE_HASH=$USE_HASH
ENV USE_DB=$USE_DB

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY=https://proxy.golang.org \
    GOARCH=amd64 \
    USER=appuser \
    UID=10001

RUN apk add --update --no-cache make git gcc build-base curl jq

#https://docs.docker.com/language/golang/build-images/

WORKDIR $GOPATH/service

COPY . .

RUN go mod download && go mod tidy && \
    go install cmd/$APP_NAME-server/main.go
RUN go build -o bin/$APP_NAME-server ./cmd/$APP_NAME-server

# STEP 2 dapp executable binary
FROM alpine:latest

# get GOPATH variable
FROM go-builder
ENV GOPATH ${GOPATH}

ARG APP_NAME=service
ARG PORT=1111
ARG NODE_ENV=production
ARG DOMAIN=localhost

# PostgreSQL settings
ARG DB_SQL_HOST=localhost
ARG DB_SQL_PORT=5432
ARG DB_SQL_USERNAME=services
ARG DB_SQL_PASSWORD=services
ARG DB_SQL_DATABASE=dapp

ARG SESSION_SECRET
ARG X_TOKEN
ARG USE_HASH
ARG USE_DB

# set environment variables
ENV APP_NAME=$APP_NAME
ENV PORT=$PORT
ENV NODE_ENV=$NODE_ENV
ENV USE_HASH=$USE_HASH
ENV USE_DB=$USE_DB
# PostgreSQL settings
ENV DB_SQL_HOST=$DB_SQL_HOST
ENV DB_SQL_PORT=$DB_SQL_PORT
ENV DB_SQL_USERNAME=$DB_SQL_USERNAME
ENV DB_SQL_PASSWORD=$DB_SQL_PASSWORD
ENV DB_SQL_DATABASE=$DB_SQL_DATABASE

ENV SESSION_SECRET=$SESSION_SECRET
ENV X_TOKEN=$X_TOKEN

WORKDIR /usr/local/bin/

# copy compiled binary and start the app
COPY --from=go-builder .$GOPATH/service/bin/$APP_NAME-server ./$APP_NAME-server
COPY --from=go-builder .$GOPATH/service/uploads ./uploads

ENTRYPOINT ./$APP_NAME-server --port=$PORT --host="0.0.0.0"
EXPOSE $PORT
