FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

RUN ALPINE_VERSION=$(cat /etc/alpine-release | cut -d'.' -f1-2) && \
    echo "https://linux-mirror.liara.ir/repository/alpine/v${ALPINE_VERSION}/main" > /etc/apk/repositories && \
    echo "https://linux-mirror.liara.ir/repository/alpine/v${ALPINE_VERSION}/community" >> /etc/apk/repositories


RUN go env -w GOPROXY=https://mirror.abrha.net/repository/go/,direct
RUN go env -w GOSUMDB=off

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY ./internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/
COPY ./internal/infrastructure/jwt/publicKey.pem /internal/infrastructure/jwt/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

FROM alpine:3.19

RUN ALPINE_VERSION=$(cat /etc/alpine-release | cut -d'.' -f1-2) && \
echo "https://linux-mirror.liara.ir/repository/alpine/v${ALPINE_VERSION}/main" > /etc/apk/repositories && \
echo "https://linux-mirror.liara.ir/repository/alpine/v${ALPINE_VERSION}/community" >> /etc/apk/repositories


RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
RUN mkdir -p /app/internal/infrastructure/jwt
COPY ./internal/infrastructure/jwt/privateKey.pem ./internal/infrastructure/jwt/
COPY ./internal/infrastructure/jwt/publicKey.pem ./internal/infrastructure/jwt/

COPY .env .

EXPOSE 8080

CMD ["./main"]