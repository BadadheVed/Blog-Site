
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o blogs-kafka ./main.go
# -------->Execute only the binary <-----------
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/blogs-kafka .
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh
EXPOSE 8080
CMD ["sh", "-c", "./blogs-kafka"]
