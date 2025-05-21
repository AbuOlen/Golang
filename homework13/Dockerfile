FROM golang:1.24-alpine

WORKDIR /work

RUN addgroup -S olena && adduser -S olena -G olena
RUN chown olena:olena /work

COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o /usr/bin/ ./cmd/server

EXPOSE 9090

USER olena

CMD ["/usr/bin/server"]