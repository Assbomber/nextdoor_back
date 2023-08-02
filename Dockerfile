FROM golang:1.20.5-alpine

ARG RAILWAY_ENVIRONMENT
ENV RAILWAY_ENVIRONMENT=$RAILWAY_ENVIRONMENT

WORKDIR /myzone

COPY . .
RUN touch .env
RUN go build -o build/myzone ./cmd/myzone

CMD ["./build/myzone"]

EXPOSE 8080