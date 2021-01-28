FROM alpine:latest

RUN addgroup api && \
  adduser -D -G api api

RUN mkdir /app

COPY ./bin/coffee-service /app/coffee-service

ENTRYPOINT [ "/app/coffee-service" ]