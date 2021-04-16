# BASE
FROM golang:1.15-alpine AS base

RUN apk add --no-cache git

WORKDIR /src

# BUILD
FROM base AS build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -o /syncodoc .

# RELEASE
FROM scratch AS release

COPY --from=build /syncodoc /syncodoc

ENTRYPOINT ["/syncodoc"]
