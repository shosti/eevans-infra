FROM golang:1.15.8 AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -tags 'osusergo netgo'

FROM gcr.io/distroless/static-debian10
WORKDIR /app
COPY --from=builder /build/cievents /app/cievents
ENTRYPOINT ["/app/cievents"]
