FROM golang:1.22-alpine AS build
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w"

# The run docker image does not need any buildtools, so we can use a plain Alpine Linux image.
FROM alpine AS run

WORKDIR /app

# Copy the binary from the "build" image.
COPY --from=build /app/backend .

# Expose SMTP send port and golang server port
EXPOSE 8080 587

# Run
CMD ["./backend"]