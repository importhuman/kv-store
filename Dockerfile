FROM golang:alpine
WORKDIR /store
COPY . .
RUN go mod download && \
    go build -o kvstore
EXPOSE 8080
CMD ["./kvstore"]