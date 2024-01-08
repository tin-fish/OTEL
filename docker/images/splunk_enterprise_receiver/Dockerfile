# Use a base image with Go and other necessary dependencies
FROM golang:alpine
RUN apk add --no-cache git
# Set the working directory
WORKDIR /otel-collector

RUN GO111MODULE=on go install go.opentelemetry.io/collector/cmd/builder@latest
COPY ./build.yaml ./
RUN builder --config=build.yaml 
COPY ./otelcol.yaml /tmp/

#CMD ["/tmp/dist/otelcol-custom --config=/tmp/otelcol.yaml"]
CMD ["/bin/sh"]
