# Stage 0: building the binary
FROM golang:alpine AS build

RUN apk update && \
    apk add curl \
            git \
            bash \
            make \
            ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN make build-main

# Stage 1: copying files
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

WORKDIR /app/

COPY --from=build /app/bin/main .
COPY --from=build /app/scripts/entrypoint.main.sh ./entrypoint.sh
COPY --from=build /app/files/service-account.json ./files/service-account.json

RUN ls -la
RUN chmod +x ./entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
