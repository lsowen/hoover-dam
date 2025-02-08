# Build the binary
FROM golang:1.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/ ./cmd/
COPY ./pkg/ ./pkg/

# Build
RUN curl --remote-name --silent https://raw.githubusercontent.com/treeverse/lakeFS/v1.33.0/api/authorization.yml && \
    go generate pkg/api/server.go
RUN CGO_ENABLED=0 go build -o hoover-dam ./cmd/hooverdam


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot@sha256:42d15c647a762d3ce3a67eab394220f5268915d6ddba9006871e16e4698c3a24
COPY --from=builder /build/hoover-dam /usr/local/bin/

USER 65532:65532

ENTRYPOINT ["/usr/local/bin/hoover-dam"]
