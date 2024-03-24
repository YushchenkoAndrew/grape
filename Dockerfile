FROM golang:1.22-alpine AS builder
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
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN make build

# Create final image
FROM alpine AS runner
WORKDIR /app

ENV GIN_MODE=release

# Copy config file & complied file
COPY configs/config.template.yaml ./configs/config.yaml 
COPY configs/database.template.yaml ./configs/database.yaml 
COPY --from=builder /app/grape .

EXPOSE 31337
CMD ["./grape"]
