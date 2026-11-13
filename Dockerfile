# Builds from source. Works with both BuildKit and the legacy builder.
FROM golang:1.23-alpine AS build
ARG VERSION=dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /lfe .

FROM scratch
COPY --from=build /lfe /lfe
ENV LFE_ROOT=/mnt LFE_LISTEN=:8080
EXPOSE 8080
ENTRYPOINT ["/lfe"]
