FROM golang:1.20.5-alpine

ARG RAILWAY_ENVIRONMENT
ENV RAILWAY_ENVIRONMENT=$RAILWAY_ENVIRONMENT

WORKDIR /myzone

COPY . .
RUN apk --update add redis 
RUN touch .env
RUN go build -o build/myzone ./cmd/myzone

CMD ["sh","-c","/usr/bin/redis-server --daemonize yes && ./build/myzone"]

EXPOSE 8080