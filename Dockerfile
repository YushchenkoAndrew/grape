FROM golang:1.17-alpine AS builder
RUN apk --no-cache add ca-certificates gcc g++ make bash git
WORKDIR /app

ENV GIN_MODE=release
ENV GO111MODULE=on

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . ./

# Update docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.7.8
RUN swag init -g ./main.go --parseDependency

RUN go build -o ./api ./main.go

# Create final image
FROM alpine AS runner
WORKDIR /

ENV GIN_MODE=release

# Copy config file & complied file
COPY .env.template .env
COPY operations.yaml .
COPY --from=builder /app/api .

EXPOSE 31337
CMD ["./api"]
