FROM golang:1.13
RUN curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
COPY ./build/linux/step-goreleaser /usr/bin/step-goreleaser
