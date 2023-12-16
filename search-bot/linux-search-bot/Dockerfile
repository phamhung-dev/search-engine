FROM golang:1.21.3

ENV http_proxy socks5://10.2.94.95:5555
ENV https_proxy socks5://10.2.94.95:5555

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

CMD [ "go", "run", "cmd/main.go" ]