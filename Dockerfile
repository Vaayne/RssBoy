FROM golang:1.16-alpine as builder
#ENV CGO_ENABLED=0
RUN apk add git make gcc libc-dev

WORKDIR /flowerss
COPY . .
RUN  make build


FROM alpine
WORKDIR /flowerss
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /flowerss/bin/ /flowerss/

CMD ["/flowerss/flowerss"]
