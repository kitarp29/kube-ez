FROM golang:alpine3.17
WORKDIR /k8s-api
COPY . .
RUN apk update && \
    apk upgrade && \
    apk --update add \
        gcc \
        g++ \
        build-base && \
    go mod tidy
CMD ["go","run","server.go"]
EXPOSE 8000