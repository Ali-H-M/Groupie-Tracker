FROM golang:1.23-alpine AS builder

LABEL version="1.0.0"
LABEL description="Go app for ascii art web. You can generate pretty ascii art through the website"
LABEL org.opencontainers.image.source="https://learn.reboot01.com/git/alimadan/groupie-tracker"

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /app ./

EXPOSE 8080

CMD [ "./main" ]