FROM golang:latest AS backend
WORKDIR /app
COPY *.go *.mod ./
RUN go mod download && go build

FROM node:latest AS frontend
WORKDIR /app
COPY frontend .
RUN npm i && npm run build

FROM scratch:latest
WORKDIR /app
COPY --from=backend /app/wake-on-lan-server-go /app/wake-on-lan-server-go
COPY --from=frontend /app/frontend /app/frontend

CMD ["wake-on-lan-server-go"]