FROM golang:1.22 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /search-index-service ./cmd/search-index-service
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /search-index-service /search-index-service
EXPOSE 8083
ENTRYPOINT ["/search-index-service"]
