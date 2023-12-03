FROM golang:1.21.4-bullseye AS builder
RUN git clone https://github.com/mdlayher/apcupsd_exporter /build
WORKDIR /build
RUN CGO_ENABLED=0 go build ./cmd/apcupsd_exporter

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /build/apcupsd_exporter /app/apcupsd_exporter

CMD ["/app/apcupsd_exporter"]
