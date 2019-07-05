# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/knative-sample/event-press
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o event-press-receive github.com/knative-sample/event-press/cmd/receive

# Copy the event-press-receive into a thin image
FROM alpine:3.7
WORKDIR /
COPY --from=builder /go/src/github.com/knative-sample/event-press/event-press-receive app/
ENTRYPOINT ["/app/event-press-receive"]