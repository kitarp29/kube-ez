FROM golang:latest
WORKDIR /k8s-api
COPY . .
RUN go mod tidy
CMD ["go","run","server.go"]
EXPOSE 8000