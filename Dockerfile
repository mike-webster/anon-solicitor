FROM golang:1.11-alpine3.8

RUN apk add --no-cache ca-certificates git make curl mysql-client gcc musl-dev

WORKDIR /anon-solicitor

ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -o anon-solicitor ./cmd/app

# For Web
EXPOSE 3001

ENTRYPOINT ["./anon-solicitor"]