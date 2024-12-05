FROM golang:1.23.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /boardwalk-mp-server

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /boardwalk-mp-server /boardwalk-mp-server

COPY *.pem ./

EXPOSE 443

USER nonroot:nonroot

ENTRYPOINT ["/boardwalk-mp-server"]