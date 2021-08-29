# ./Dockerfile

FROM golang:1.17-alpine AS build

RUN apk --no-cache add build-base=0.5-r2

WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY ./src .

RUN CGO_ENABLED=1 go build -o /bin/app .

FROM alpine:3

COPY --from=build /bin/app /bin/app

ENV SUBSTITUTIONS_UPDATECRON="*/10 6-23 * * *" MOODLE_UPDATECRON="0 6-23 * * *" DATABASE_URI="db.sqlite" DATABASE_TYPE="SQLITE" SIGNAL_CLI_GRPC_API_URL="localhost:9000" DATABASE_AUTOMIGRATE=1

ENTRYPOINT ["/bin/app"]