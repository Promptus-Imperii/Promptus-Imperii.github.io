# Base images are pinned by digest; Dependabot updates the digests.
FROM golang:1.27-alpine@sha256:eb61c3873c97416fff0116553fff65d67a3549f72e11d9758849a689cdd8b523 AS build
WORKDIR /src
COPY go.mod go.sum ./
# The cache mounts keep downloaded modules and compiled packages between builds.
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY *.go ./
COPY templates ./templates
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /backend .

# distroless/static has CA certificates (needed for SMTP over TLS) and runs as non-root.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /backend /backend
EXPOSE 3000
ENTRYPOINT ["/backend"]
