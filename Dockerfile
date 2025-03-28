# Stage 1: Build Go API
FROM golang:1.18-alpine AS api-builder

# Install build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy Go project files
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download

# Copy API source code
COPY apps/api/ .

# Build Go binary
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/api

# Stage 2: Build Node.js Web App
FROM node:18-alpine AS web-builder
WORKDIR /app

# Copy web app files
COPY apps/web/package*.json ./
RUN npm install

# Copy web app source
COPY apps/web/ .

# Build web app
RUN npm run build

# Stage 3: Final image with both services
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache sqlite-libs nodejs npm

WORKDIR /app

# Copy Go API binary
COPY --from=api-builder /app/api ./

# Copy built Node.js web app
COPY --from=web-builder /app ./web

# Expose ports for both services
EXPOSE 8080 3000

# Create an entrypoint script to run both services
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set environment variables
ENV PORT=8080 \
    SQLITE_PATH=/data/sqlite/users.db \
    MONGO_URI=mongodb://root:strongpassword@mongodb:27017 \
    JWT_SECRET=supersecretkey \
    NODE_ENV=production

# Use the entrypoint script
ENTRYPOINT ["/entrypoint.sh"]