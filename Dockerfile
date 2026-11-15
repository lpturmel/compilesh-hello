FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /hello .

# distroless "nonroot" runs as the numeric user 65532 and has no shell.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /hello /hello
ENTRYPOINT ["/hello"]
