FROM golang:alpine AS build
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
# -ldflags "-s -w"
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine AS run

WORKDIR /app
RUN apk add --no-cache msmtp postfix

COPY --from=build /app/backend .

# Copy msmtp configuration file
COPY docker/msmtp.conf /etc/msmtprc
COPY docker/main.cf /etc/postfix/main.cf
COPY docker/start.sh .
RUN touch /etc/postfix/recipient_validation && postmap /etc/postfix/recipient_validation

RUN chmod +x start.sh
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080
# Expose SMTP port
EXPOSE 25 587

# Run
# Start Postfix service and sm
CMD ["./start.sh"]