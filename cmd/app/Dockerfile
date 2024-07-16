FROM  golang:alpine AS base
WORKDIR /go/tranched
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/app

FROM alpine:latest
COPY --from=base /app /app
ENTRYPOINT [ "/app" ]