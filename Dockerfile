# Build Stage
FROM golang:1.21-alpine3.18 as BUILDER
LABEL authors="daovudat"
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.18
WORKDIR /app
COPY --from=BUILDER /app/main .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]