# ./Dockerfile

FROM golang:1.17-alpine AS builder

# Move to working directory (/build).
RUN apk --no-cache add ca-certificates=20161130-r0 

WORKDIR /build

# Copy and download dependency using go mod.
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy the code into the container.
COPY ./src .

# Set necessary environment variables needed for our image 
# and build the API server.
RUN CGO_ENABLED=1 go build -o app .

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app /

# Necessary environment variables
ENV SUBSTITUTIONS_UPDATECRON="*/10 6-23 * * *" MOODLE_UPDATECRON="0 6-23 * * *" DATABASE_URI="db.sqlite" DATABASE_TYPE="SQLITE" SIGNAL_CLI_GRPC_API_URL="localhost:9000" DATABASE_AUTOMIGRATE=1

# Command to run when starting the container.
ENTRYPOINT ["/app"]