# Build app
FROM golang:1.23 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -C cmd/api -o /app/go-api-products


# Run app
FROM cgr.dev/chainguard/static:latest AS release

WORKDIR /app

COPY --from=build /app/go-api-products ./

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT [ "/app/go-api-products" ]