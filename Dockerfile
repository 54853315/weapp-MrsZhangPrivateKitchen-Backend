FROM golang:1.13.5-alpine3.10 AS builder

WORKDIR /build
RUN adduser -u 10001 -D app-runner
#RUN export DOCKER_HOST_IP=$(route -n | awk '/UG[ \t]/{print $2}')

ENV GOPROXY https://goproxy.cn
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o FoodBackend .

FROM alpine:3.10 AS final

WORKDIR /app
COPY ./ /app/
COPY --from=builder /build/FoodBackend /app/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER app-runner
ENTRYPOINT ["/app/FoodBackend"]