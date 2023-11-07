FROM golang:1.21 AS interim
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY narrative.yaml .
RUN go get ./...
RUN go build -o otlpumper
RUN if [ -f /app/otlpumper ]; then echo "otlpumper file exists"; else echo "otlpumper file does not exist"; fi
FROM busybox
WORKDIR /app
COPY --from=interim /app/otlpumper /app/otlpumper
COPY narrative.yaml .
CMD ["./otlpumper"]
