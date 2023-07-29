FROM golang:1.20.5-alpine
WORKDIR /myzone
COPY . .
RUN go build -o build/myzone .
CMD ["./build/myzone"]
EXPOSE 8080