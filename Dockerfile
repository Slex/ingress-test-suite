ARG VERSION=dev
FROM golang:1.25-alpine AS builder
ARG VERSION
WORKDIR /build
COPY . .
RUN apk add make
RUN make build

FROM alpine:3.21.3
ENV MAIN_SERVER_PORT=8080
ENV METRICS_SERVER_PORT=9091
WORKDIR /app
COPY --from=builder /build/bin/ingress-test-suite /app/ingress-test-suite
RUN chmod +x /app/ingress-test-suite
CMD ["/app/ingress-test-suite"]
