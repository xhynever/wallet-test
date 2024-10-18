# FROM golang:alpine as builder
# WORKDIR /app
# COPY . .
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./main ./main.go 


# Run stage
FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY ./main .
COPY .env .
# COPY start.sh .
# COPY wait-for.sh .
# COPY db/migration ./db/migration
EXPOSE 8080 
CMD ["./main"]
