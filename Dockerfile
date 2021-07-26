FROM golang:1.16.6-alpine3.14 AS builder
WORKDIR /oem

# Download dependencies first, so they are not pulled every
# time any minor code changes are made
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build ./cmd/online-election-manager

FROM alpine:3.14
WORKDIR /oem
COPY --from=builder /oem/online-election-manager ./
ENTRYPOINT ["./online-election-manager"]
