FROM golang:latest AS backend
WORKDIR /app
COPY *.go *.mod ./
RUN go mod download
RUN go build -a -tags netgo -ldflags '-w'

FROM node:latest AS frontend
WORKDIR /app
COPY frontend .
RUN npm ci && npm run build

FROM alpine:latest
LABEL org.opencontainers.image.source "https://github.com/Tedyst/wake-on-lan-server-go"

WORKDIR /app
COPY --from=backend /app/wake-on-lan-server-go /app/wake-on-lan-server-go
COPY --from=frontend /app/build /app/frontend/build

ENTRYPOINT ["/app/wake-on-lan-server-go"]