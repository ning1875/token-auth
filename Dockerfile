FROM golang:latest as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./




RUN go mod download
COPY . .

RUN CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build -o token-auth auth.go


#FROM yauritux/busybox-curl  as runner
FROM alpine as runner
COPY --from=builder /app/token-auth .
ENTRYPOINT [ "./token-auth" ]

