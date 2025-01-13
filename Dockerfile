FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /build

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd


FROM alpine:latest

RUN mkdir -p /home

COPY --from=builder /build/app /home/app

EXPOSE 9010
EXPOSE 9050

CMD ["./home/app"]


