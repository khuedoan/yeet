FROM docker.io/golang:1.23-alpine AS builder

WORKDIR /src

COPY . .

RUN go build -o /bin/worker ./worker
RUN go build -o /bin/server ./server

FROM scratch

COPY --from=builder /bin/worker /bin/worker
COPY --from=builder /bin/server /bin/server

CMD ["/bin/worker"]
