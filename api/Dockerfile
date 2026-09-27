### Step 1: Build stage
FROM golang:1.25.5-alpine3.22 as builder
ARG CGO_ENABLED=0
ARG GOOS=linux

# Create a non-root user for running the application
#RUN useradd -u 1001 nonroot

WORKDIR /app

# This will download all certificates (ca-certificates) and builds it in a
# single file under /etc/ssl/certs/ca-certificates.crt (update-ca-certificates)
# I also add git so that we can download with `go mod download` and
# tzdata to configure timezone in final image
RUN apk --update add --no-cache ca-certificates openssl git tzdata && \
    update-ca-certificates

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./

RUN go mod download

# Copy the application source code and build the binary
COPY . .

# Compile the application with static linking
RUN go build -o kfamily ./cmd/api

###
## Step 2: Runtime stage
FROM scratch
LABEL maintainer="M Kumaran"

# Copy /etc/passwd for non-root user information
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# This line will copy all certificates to final image
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the binary from the build stage to the final image
COPY --from=builder /app/kfamily /

# Use the non-root user
#USER nonroot

EXPOSE 8080
# Set the entry point for the container
ENTRYPOINT ["/kfamily"]
