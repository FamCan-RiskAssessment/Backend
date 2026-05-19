FROM golang:1.23.4-alpine AS builder

RUN ver="v$(cut -d. -f1,2 /etc/alpine-release)" && \
    printf '%s\n' \
      "https://mirror.arvancloud.ir/alpine/${ver}/main" \
      "https://mirror.arvancloud.ir/alpine/${ver}/community" \
      > /etc/apk/repositories

RUN apk update && apk add --no-cache gcc musl-dev git && rm -rf /var/cache/apk/*

WORKDIR /app

RUN go env -w GOPROXY=https://package-mirror.liara.ir/repository/go/
RUN go env -w GOSUMDB=off

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

FROM alpine:3.19

RUN ver="v$(cut -d. -f1,2 /etc/alpine-release)" && \
    printf '%s\n' \
      "https://mirror.arvancloud.ir/alpine/${ver}/main" \
      "https://mirror.arvancloud.ir/alpine/${ver}/community" \
      > /etc/apk/repositories

RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
