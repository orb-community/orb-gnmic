FROM golang:1.20 AS builder
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/orb-community/orb-gnmic
COPY . .
RUN apt update -y \
    && apt install make build-essential dpkg-dev curl ca-certificates libasan6 wget gzip -y
RUN make getgnmic \
    && rm -rf build \
    && CGO_ENABLED=1 make build \
    && mv build/orb-gnmic /exe

FROM debian:bullseye-slim

RUN mkdir /opt/orb-gnmic/

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
