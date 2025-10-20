FROM golang:1.24-alpine AS backend-deps

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine AS backend-build

WORKDIR /app

COPY --from=backend-deps /go/pkg /go/pkg
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
  CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o cmd ./registry/cmd/main.go

FROM scratch

COPY --from=backend-build /app/cmd /cmd

EXPOSE 8080

CMD ["./cmd"]