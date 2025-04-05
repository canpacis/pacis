FROM golang:1.24.1-alpine AS builder

RUN apk add --no-cache git
  
WORKDIR /app
  
COPY . .
  
RUN go build -o pacis-docs ./docs

FROM alpine:latest

RUN adduser -D -g '' appuser

COPY --from=builder /app/pacis-docs /usr/local/bin/pacis-docs
COPY --from=builder /app/docs/app/markup /usr/local/bin/docs/app/markup

COPY ./docs /app/docs

USER appuser

EXPOSE 8080
ENTRYPOINT ["pacis-docs"]