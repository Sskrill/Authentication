FROM golang:alpine AS builder
WORKDIR /app
RUN apk --no-cache add bash
COPY ["go.mod","go.sum", "./"]
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
RUN go build -o ./bin/app cmd/main.go

FROM alpine
COPY --from=builder app/bin/app /
COPY .env /
CMD ["/app"]