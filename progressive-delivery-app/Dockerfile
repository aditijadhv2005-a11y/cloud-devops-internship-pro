FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .
FROM gcr.io/distroless/static:nonroot
COPY --from=build /out/server /server
USER 65532:65532
EXPOSE 8080
ENV APP_VERSION=container
ENTRYPOINT ["/server"]
