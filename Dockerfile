# Step 1: Use official Go image as the build environment
FROM golang:1.23-alpine AS build

# Step 2: Set the current working directory inside the container
WORKDIR /app

# Step 3: Copy the go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Step 4: Download dependencies
RUN go mod download

# Step 5: Copy the rest of the application code
COPY . .

# Step 6: Build the Go app
RUN go build -o main .

# Step 7: Use a minimal image for running the app
FROM alpine:latest

# Step 8: Install certificates (required for database connections)
RUN apk --no-cache add ca-certificates

# Step 9: Set working directory
WORKDIR /root/

# Step 10: Copy the binary from the build stage
COPY --from=build /app/main .

# Step 11: Copy the .env file (if any) to the container (optional)
COPY .env .env

# Step 12: Expose the port for the app
EXPOSE 8080

# Step 13: Command to run the executable
CMD ["./main"]
