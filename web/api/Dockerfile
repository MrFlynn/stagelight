# Build container.
FROM golang:latest as build

WORKDIR /go/src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o api .

# Actually runnable container.
FROM scratch

VOLUME [ "/data" ]

EXPOSE 8000/tcp

WORKDIR /
COPY --from=build /go/src/api /
CMD ["/api", "--address=0.0.0.0", "--port=8000", "--database=/data/example.db"]