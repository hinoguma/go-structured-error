FROM golang:1.22-alpine

RUN apk add --no-cache git curl

ENV GO111MODULE=on
ENV GOPATH=/gopath

# Install golangci-lint (adjust version as needed)
ENV GOLANGCI_LINT_VERSION=v2.5.0
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh \
    | sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
