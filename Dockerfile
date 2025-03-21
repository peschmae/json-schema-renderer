FROM golang:alpine as build

ENV CGO_ENABLED=0

WORKDIR /root
COPY . /root
RUN go build -ldflags="-w -s" .

FROM alpine:latest

RUN addgroup -S appuser \
    && adduser -S appuser -G appuser \
    && mkdir /app \
    && chown -R appuser:appuser /app

COPY --from=build /root/json-schema-renderer /usr/bin/json-schema-renderer

USER appuser
WORKDIR /app

CMD ["json-schema-renderer"]
