FROM golang:1.22.6-alpine3.20 as builder

COPY . /src/
WORKDIR /src

RUN apk add --no-cache make=4.4.1-r2 && make -j "$(nproc)"

FROM alpine:3.20

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

COPY --from=builder /src/dist/ponyhug /app/ponyhug

USER 1000
ENTRYPOINT [ "/app/ponyhug" ]
