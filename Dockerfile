FROM golang:alpine AS BUILDER
WORKDIR /usr/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY cmd/ ./cmd
COPY internal/ ./internal
WORKDIR /usr/src/app/cmd
RUN go build -o ./

FROM alpine:latest 
RUN apk --no-cache add ca-certificates
WORKDIR /usr/src/app
EXPOSE 8080
COPY cmd/migrations/ ./migrations
COPY --from=BUILDER /usr/src/app/cmd/ ./
ENTRYPOINT [ "./cmd" ]