FROM golang:alpine AS build
# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build

FROM alpine AS run

WORKDIR /app
# Install msmtp and postfix: msmtp is the CLI mail client, postfix is the smtp server
RUN apk add --no-cache msmtp postfix

COPY --from=build /app/backend .

# Copy msmtp and postfix configuration file
COPY docker/msmtp.conf /etc/msmtprc
COPY docker/main.cf /etc/postfix/main.cf

# Copy start.sh
COPY docker/start.sh .

RUN mkdir /app/inschrijvingen && chmod +x start.sh

# Expose SMTP port and golang server port
EXPOSE 25 8080

# Run
# Start Postfix service and sm
CMD ["./start.sh"]