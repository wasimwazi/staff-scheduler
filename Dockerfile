FROM golang:latest as backend
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o scheduler ./cmd

FROM alpine
WORKDIR /app
COPY --from=backend /app/scheduler .
CMD ./scheduler