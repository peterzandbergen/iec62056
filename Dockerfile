FROM golang:1.11.2 AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -v cmd/emlog/*.go

FROM scratch

WORKDIR /app

COPY --from=builder /build/emlog /app/emlog

ENTRYPOINT [ "/app/emlog" ]

