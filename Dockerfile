# Build

FROM golang:1.26.3-alpine3.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./executable --ldflags "-s -w" ./cmd/main.go

# Release

FROM gcr.io/distroless/static-debian13

WORKDIR /app

COPY --from=build /app/executable ./

EXPOSE 8080
ENTRYPOINT [ "/app/executable" ]
