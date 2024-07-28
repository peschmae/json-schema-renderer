FROM golang:alpine as build

WORKDIR /root
COPY . /root
RUN go build .

FROM alpine:latest

RUN addgroup -S appuser \
    && adduser -S appuser -G appuser \
    && mkdir /app \
    && chown -R appuser:appuser /app

COPY --from=build /root/json-schema-asciidoc /usr/bin/json-schema-asciidoc

USER appuser
WORKDIR /app

CMD ["json-schema-asciidoc"]
