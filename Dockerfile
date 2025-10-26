FROM golang:1.25-alpine AS builder

WORKDIR /appGo

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

RUN go build -o my-app ./

FROM alpine:latest

WORKDIR /app

COPY --from=builder /appGo/my-app /app/my-app

EXPOSE 8080

ENTRYPOINT [ "./my-app" ]