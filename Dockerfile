FROM golang:1.23.3-alpine AS BUILDER

WORKDIR /go/src/github.com/zaidsasa/xbankapi

RUN apk add make curl

COPY . .

RUN mkdir -p /tools/ && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
    -o /tools/migrate.linux-amd64.tar.gz && \
    tar xzvf /tools/migrate.linux-amd64.tar.gz -C /tools

RUN make build

FROM alpine:3.20

LABEL org.opencontainers.image.source="https://github.com/zaidsasa/xbankapi"

WORKDIR /app

COPY --from=BUILDER /tools/migrate /app/migrate
COPY --from=BUILDER /go/src/github.com/zaidsasa/xbankapi/entrypoint.sh /app/entrypoint.sh
COPY --from=BUILDER /go/src/github.com/zaidsasa/xbankapi/out/xbankapi /app/xbankapi
COPY --from=BUILDER /go/src/github.com/zaidsasa/xbankapi/db/migrations/ /app/migrations/

ENV DATABASE_URL=

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chown -R appuser:appgroup /app 

USER appuser

ENTRYPOINT ["/app/entrypoint.sh"]

EXPOSE  3000
