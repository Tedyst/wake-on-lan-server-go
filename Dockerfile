FROM golang:latest AS backend
WORKDIR /app
COPY *.go *.mod ./
RUN go mod download && go build -a -tags netgo -ldflags '-w'

FROM node:latest AS frontend
WORKDIR /app
COPY frontend .
RUN npm ci && npm run build

FROM alpine:latest
WORKDIR /app
COPY --from=backend /app/wake-on-lan-server-go /app/wake-on-lan-server-go
COPY --from=frontend /app/build /app/frontend/build

CMD ["/app/wake-on-lan-server-go"]