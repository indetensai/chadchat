FROM golang:alpine AS builder
WORKDIR /usr/src/chadchat
COPY go.mod go.sum cmd/main.go cmd/generate-keys.sh ./
COPY internal ./internal/
RUN go mod download &&\
    go build -o ./server
COPY cmd/*.pem ./
RUN apk add openssl &&\
    chmod +x /usr/src/chadchat/generate-keys.sh &&\
    ./generate-keys.sh

FROM alpine:latest AS runner
WORKDIR /usr/bin/chadchat
EXPOSE 8080
COPY cmd/migrations ./migrations/
COPY --from=builder /usr/src/chadchat/server /usr/src/chadchat/*.pem ./
RUN apk --no-cache add ca-certificates
ENTRYPOINT [ "/usr/bin/chadchat/server" ]