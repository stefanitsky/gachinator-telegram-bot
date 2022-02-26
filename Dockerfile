# Build stage 
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /bot

# Deploy stage
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bot /bot

USER nonroot:nonroot

ENTRYPOINT ["/bot"]
