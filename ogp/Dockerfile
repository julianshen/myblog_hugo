FROM golang:1.16.6-alpine

RUN apk add --update \
  --repository http://dl-3.alpinelinux.org/alpine/edge/testing \
  vips-dev libc-dev gcc \
  && rm -rf /var/cache/apk/*
COPY . /app
WORKDIR /app
RUN go build .
ENTRYPOINT [ "/app/ogp", "-b", "/blog", "-f" ]
