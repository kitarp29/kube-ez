# Build
FROM golang:alpine AS build

WORKDIR /k8-api
COPY . .
RUN go mod download && \
    GO111MODULE=on CGO_ENABLED=0 go build -ldflags "-s -w"

# Deploy
FROM alpine
COPY --from=build /k8-api/k8-api .

ENTRYPOINT ["./k8-api"]
# CMD WALA PART?
EXPOSE 8000